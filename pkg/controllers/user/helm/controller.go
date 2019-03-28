package helm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"k8s.io/api/core/v1"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/rancher/rancher/pkg/controllers/management/compose/common"
	"github.com/rancher/rancher/pkg/ref"
	"github.com/rancher/rancher/pkg/systemaccount"
	corev1 "github.com/rancher/types/apis/core/v1"
	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	helmcommon "github.com/rancher/rancher/pkg/controllers/user/helm/common"
)

const (
	helmTokenPrefix           = "helm-token-"
	description               = "token for helm chart deployment"
	AppIDsLabel               = "cattle.io/appIds"
	creatorIDAnn              = "field.cattle.io/creatorId"
	MultiClusterAppIDSelector = "mcapp"
	projectIDFieldLabel       = "field.cattle.io/projectId"
)

func Register(ctx context.Context, user *config.UserContext, kubeConfigGetter common.KubeConfigGetter) {
	appClient := user.Management.Project.Apps("")
	stackLifecycle := &Lifecycle{
		KubeConfigGetter:         kubeConfigGetter,
		SystemAccountManager:     systemaccount.NewManager(user.Management),
		TokenClient:              user.Management.Management.Tokens(""),
		UserClient:               user.Management.Management.Users(""),
		UserManager:              user.Management.UserManager,
		K8sClient:                user.K8sClient,
		TemplateVersionClient:    user.Management.Management.CatalogTemplateVersions(""),
		TemplateClient:           user.Management.Management.CatalogTemplates(""),
		CatalogLister:            user.Management.Management.Catalogs("").Controller().Lister(),
		ClusterCatalogLister:     user.Management.Management.ClusterCatalogs("").Controller().Lister(),
		ProjectCatalogLister:     user.Management.Management.ProjectCatalogs("").Controller().Lister(),
		ListenConfigClient:       user.Management.Management.ListenConfigs(""),
		ClusterName:              user.ClusterName,
		AppRevisionGetter:        user.Management.Project,
		AppGetter:                user.Management.Project,
		NsLister:                 user.Core.Namespaces("").Controller().Lister(),
		NsClient:                 user.Core.Namespaces(""),
		istioClusterGraphClient:  user.Management.Management.IstioClusterMonitorGraphs(metav1.NamespaceAll),
		istioProjectGraphClient:  user.Management.Management.IstioProjectMonitorGraphs(metav1.NamespaceAll),
		istioMonitorMetricClient: user.Management.Management.MonitorMetrics(metav1.NamespaceAll),
		projectGetter:            user.Management.Management.Projects(user.ClusterName),
	}
	appClient.AddClusterScopedLifecycle(ctx, "helm-controller", user.ClusterName, stackLifecycle)

	StartStateCalculator(ctx, user)
}

type Lifecycle struct {
	KubeConfigGetter         common.KubeConfigGetter
	SystemAccountManager     *systemaccount.Manager
	UserManager              user.Manager
	TokenClient              mgmtv3.TokenInterface
	UserClient               mgmtv3.UserInterface
	TemplateVersionClient    mgmtv3.CatalogTemplateVersionInterface
	TemplateClient           mgmtv3.CatalogTemplateInterface
	CatalogLister            mgmtv3.CatalogLister
	ClusterCatalogLister     mgmtv3.ClusterCatalogLister
	ProjectCatalogLister     mgmtv3.ProjectCatalogLister
	K8sClient                kubernetes.Interface
	ListenConfigClient       mgmtv3.ListenConfigInterface
	ClusterName              string
	AppRevisionGetter        v3.AppRevisionsGetter
	AppGetter                v3.AppsGetter
	NsLister                 corev1.NamespaceLister
	NsClient                 corev1.NamespaceInterface
	istioClusterGraphClient  mgmtv3.IstioClusterMonitorGraphInterface
	istioProjectGraphClient  mgmtv3.IstioProjectMonitorGraphInterface
	istioMonitorMetricClient mgmtv3.MonitorMetricInterface
	projectGetter            mgmtv3.ProjectInterface
}

func (l *Lifecycle) Create(obj *v3.App) (runtime.Object, error) {
	v3.AppConditionMigrated.True(obj)
	return obj, nil
}

func (l *Lifecycle) Updated(obj *v3.App) (runtime.Object, error) {
	if obj.Spec.ExternalID == "" && len(obj.Spec.Files) == 0 {
		return obj, nil
	}
	// always refresh app to avoid updating app twice
	_, projectName := ref.Parse(obj.Spec.ProjectName)
	var err error
	obj, err = l.AppGetter.Apps(projectName).Get(obj.Name, metav1.GetOptions{})
	if err != nil {
		return obj, err
	}
	// if app was created before 2.1, run migrate to install a no-op helm release
	newObj, err := v3.AppConditionMigrated.Once(obj, func() (runtime.Object, error) {
		return l.DeployApp(obj)
	})
	if err != nil {
		return obj, err
	}
	obj = newObj.(*v3.App)
	appRevisionClient := l.AppRevisionGetter.AppRevisions(projectName)
	if obj.Spec.AppRevisionName != "" {
		currentRevision, err := appRevisionClient.Get(obj.Spec.AppRevisionName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if obj.Spec.ExternalID != "" {
			if currentRevision.Status.ExternalID == obj.Spec.ExternalID && reflect.DeepEqual(currentRevision.Status.Answers, obj.Spec.Answers) && reflect.DeepEqual(currentRevision.Status.ValuesYaml, obj.Spec.ValuesYaml) {
				if !v3.AppConditionForceUpgrade.IsTrue(obj) {
					v3.AppConditionForceUpgrade.True(obj)
				}
				return obj, nil
			}
		}
		if obj.Status.AppliedFiles != nil {
			if reflect.DeepEqual(obj.Status.AppliedFiles, obj.Spec.Files) && reflect.DeepEqual(currentRevision.Status.Answers, obj.Spec.Answers) && reflect.DeepEqual(currentRevision.Status.ValuesYaml, obj.Spec.ValuesYaml) {
				if !v3.AppConditionForceUpgrade.IsTrue(obj) {
					v3.AppConditionForceUpgrade.True(obj)
				}
				return obj, nil
			}
		}
	}
	created := false
	if obj.Spec.MultiClusterAppName != "" {
		if _, err := l.NsLister.Get("", obj.Spec.TargetNamespace); err != nil && !errors.IsNotFound(err) {
			return obj, err
		} else if err != nil && errors.IsNotFound(err) {
			n := v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: obj.Spec.TargetNamespace,
					Labels: map[string]string{
						projectIDFieldLabel:       obj.Namespace,
						MultiClusterAppIDSelector: obj.Spec.MultiClusterAppName},
					Annotations: map[string]string{
						projectIDFieldLabel: fmt.Sprintf("%s:%s", l.ClusterName, obj.Namespace),
						creatorIDAnn:        obj.Annotations[creatorIDAnn],
						AppIDsLabel:         obj.Name},
				},
			}
			_, err := l.NsClient.Create(&n)
			if err != nil && !errors.IsAlreadyExists(err) {
				return obj, err
			} else if err == nil {
				created = true
			}
		}
	}
	result, err := l.DeployApp(obj)
	if err != nil {
		return result, err
	}
	if !v3.AppConditionForceUpgrade.IsTrue(obj) {
		v3.AppConditionForceUpgrade.True(obj)
	}
	ns, err := l.NsClient.Get(obj.Spec.TargetNamespace, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return result, err
		}
		if obj.Spec.MultiClusterAppName != "" {
			return obj, fmt.Errorf("namespace not found after app creation %s", obj.Name)
		}
		return result, nil
	}
	if created {
		return result, nil
	}
	if ns.Annotations[AppIDsLabel] == "" {
		ns.Annotations[AppIDsLabel] = obj.Name
	} else {
		parts := strings.Split(ns.Annotations[AppIDsLabel], ",")
		appIDs := map[string]struct{}{}
		for _, part := range parts {
			appIDs[part] = struct{}{}
		}
		appIDs[obj.Name] = struct{}{}
		appIDList := []string{}
		for k := range appIDs {
			appIDList = append(appIDList, k)
		}
		sort.Strings(appIDList)
		ns.Annotations[AppIDsLabel] = strings.Join(appIDList, ",")
	}
	if _, err := l.NsClient.Update(ns); err != nil {
		return result, err
	}

	l.deployIstioExpression(obj)

	return result, nil
}

func (l *Lifecycle) deployIstioExpression(obj *v3.App) {
	values, err := url.Parse(obj.Spec.ExternalID)
	if err != nil {
		logrus.Errorf("check catalog type failed: %s", err.Error())
	}

	catalogWithNamespace := values.Query().Get("catalog")
	split := strings.SplitN(catalogWithNamespace, "/", 2)
	catalog := split[len(split)-1]

	template := values.Query().Get("template")

	if catalog == "library" && template == "rancher-istio" {
		helmcommon.IsMetricExpressionDeployed(l.ClusterName, l.projectGetter, l.istioClusterGraphClient, l.istioProjectGraphClient, l.istioMonitorMetricClient, obj)
	}
}

func (l *Lifecycle) DeployApp(obj *v3.App) (*v3.App, error) {
	obj = obj.DeepCopy()
	var err error
	if !v3.AppConditionInstalled.IsUnknown(obj) {
		v3.AppConditionInstalled.Unknown(obj)
		obj, err = l.AppGetter.Apps("").Update(obj)
		if err != nil {
			return obj, err
		}
	}
	newObj, err := v3.AppConditionInstalled.Do(obj, func() (runtime.Object, error) {
		template, notes, appDir, tempDir, err := l.generateTemplates(obj)
		defer os.RemoveAll(tempDir)
		if err != nil {
			return obj, err
		}
		if err := l.Run(obj, template, appDir, notes); err != nil {
			return obj, err
		}
		return obj, nil
	})
	return newObj.(*v3.App), err
}

func (l *Lifecycle) Remove(obj *v3.App) (runtime.Object, error) {
	_, projectName := ref.Parse(obj.Spec.ProjectName)
	appRevisionClient := l.AppRevisionGetter.AppRevisions(projectName)
	revisions, err := appRevisionClient.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", appLabel, obj.Name),
	})
	if err != nil {
		return obj, err
	}
	for _, revision := range revisions.Items {
		if err := appRevisionClient.Delete(revision.Name, &metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return obj, err
		}
	}
	tempDir, err := ioutil.TempDir("", "helm-")
	if err != nil {
		return obj, err
	}
	defer os.RemoveAll(tempDir)
	kubeConfigPath, err := l.writeKubeConfig(obj, tempDir, true)
	if err != nil {
		return obj, err
	}
	// try three times and succeed
	start := time.Second * 1
	for i := 0; i < 3; i++ {
		if err = helmDelete(kubeConfigPath, obj); err == nil {
			break
		}
		logrus.Warn(err)
		time.Sleep(start)
		start *= 2
	}
	ns, err := l.NsClient.Get(obj.Spec.TargetNamespace, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return obj, err
	} else if errors.IsNotFound(err) {
		return obj, nil
	}
	if val, ok := ns.Labels[MultiClusterAppIDSelector]; ok && val == obj.Spec.MultiClusterAppName {
		err := l.NsClient.Delete(ns.Name, &metav1.DeleteOptions{})
		if errors.IsNotFound(err) {
			return obj, nil
		}
		return obj, err
	}
	appIds := strings.Split(ns.Annotations[AppIDsLabel], ",")
	appAnno := ""
	for _, appID := range appIds {
		if appID == obj.Name {
			continue
		}
		appAnno += appID + ","
	}
	if appAnno == "" {
		delete(ns.Annotations, AppIDsLabel)
	} else {
		appAnno = strings.TrimSuffix(appAnno, ",")
		ns.Annotations[AppIDsLabel] = appAnno
	}
	if _, err := l.NsClient.Update(ns); err != nil {
		return obj, err
	}

	l.withdrawIstioExpression(obj)

	return obj, nil
}

func (l *Lifecycle) withdrawIstioExpression(obj *v3.App) {
	values, err := url.Parse(obj.Spec.ExternalID)
	if err != nil {
		logrus.Errorf("check catalog type failed: %s", err.Error())
	}

	catalogWithNamespace := values.Query().Get("catalog")
	split := strings.SplitN(catalogWithNamespace, "/", 2)
	catalog := split[len(split)-1]

	template := values.Query().Get("template")

	if catalog == "library" && template == "rancher-istio" {
		helmcommon.IsMetricExpressionWithdraw(l.ClusterName, l.projectGetter, l.istioClusterGraphClient, l.istioProjectGraphClient, l.istioMonitorMetricClient, obj)
	}
}

func (l *Lifecycle) Run(obj *v3.App, template, templateDir, notes string) error {
	tempDir, err := ioutil.TempDir("", "kubeconfig-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	kubeConfigPath, err := l.writeKubeConfig(obj, tempDir, false)
	if err != nil {
		return err
	}
	if err := helmInstall(templateDir, kubeConfigPath, obj); err != nil {
		return err
	}
	return l.createAppRevision(obj, template, notes, false)
}

func (l *Lifecycle) createAppRevision(obj *v3.App, template, notes string, failed bool) error {
	_, projectName := ref.Parse(obj.Spec.ProjectName)
	appRevisionClient := l.AppRevisionGetter.AppRevisions(projectName)
	release := &v3.AppRevision{}
	release.GenerateName = "apprevision-"
	release.Labels = map[string]string{
		appLabel: obj.Name,
	}
	if failed {
		release.Labels[failedLabel] = "true"
	}
	release.Status.Answers = obj.Spec.Answers
	release.Status.ProjectName = projectName
	release.Status.ExternalID = obj.Spec.ExternalID
	release.Status.ValuesYaml = obj.Spec.ValuesYaml
	release.Status.Files = obj.Spec.Files

	digest := sha256.New()
	digest.Write([]byte(template))
	tag := hex.EncodeToString(digest.Sum(nil))
	release.Status.Digest = tag

	createdRevision, err := appRevisionClient.Create(release)
	if err != nil {
		return err
	}
	obj.Spec.AppRevisionName = createdRevision.Name
	obj.Status.Notes = notes
	if obj.Spec.Files != nil {
		obj.Status.AppliedFiles = obj.Spec.Files
	}
	return err
}

func (l *Lifecycle) writeKubeConfig(obj *v3.App, tempDir string, remove bool) (string, error) {
	var token string

	userID := obj.Annotations["field.cattle.io/creatorId"]
	user, err := l.UserClient.Get(userID, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return "", err
	} else if errors.IsNotFound(err) && remove {
		token, err = l.SystemAccountManager.GetOrCreateProjectSystemToken(obj.Namespace)
	} else if err == nil {
		token, err = l.UserManager.EnsureToken(helmTokenPrefix+user.Name, description, user.Name)
	}
	if err != nil {
		return "", err
	}

	kubeConfig := l.KubeConfigGetter.KubeConfig(l.ClusterName, token)
	for k := range kubeConfig.Clusters {
		kubeConfig.Clusters[k].InsecureSkipTLSVerify = true
	}
	if err := os.MkdirAll(filepath.Join(tempDir, obj.Namespace), 0755); err != nil {
		return "", err
	}
	kubeConfigPath := filepath.Join(tempDir, obj.Namespace, ".kubeconfig")
	if err := clientcmd.WriteToFile(*kubeConfig, kubeConfigPath); err != nil {
		return "", err
	}
	return kubeConfigPath, nil
}
