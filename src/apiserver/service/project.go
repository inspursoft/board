package service

import (
	"errors"
	"fmt"

	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"

	"github.com/astaxie/beego/logs"

	//modelK8s "k8s.io/client-go/pkg/api/v1"

	"git/inspursoft/board/src/common/k8sassist"
)

var repoServeURL = utils.GetConfig("REPO_SERVE_URL")

const (
	k8sAPIversion1 = "v1"
	adminUserID    = 1
	adminUserName  = "boardadmin"
	projectPrivate = 0
	kubeNamespace  = "kube-system"
	istioNamespace = "istio-system"
	istioLabel     = "istio-injection"
)

var undeletableNamespaces = []string{kubeNamespace, istioNamespace, "kube-node-lease", "kube-public", "default", "library", "kubeedge", "cadvisor"}

func CreateProject(project model.Project) (bool, error) {

	projectID, err := dao.AddProject(project)
	if err != nil {
		return false, err
	}

	projectMember := model.ProjectMember{
		ProjectID: projectID,
		UserID:    int64(project.OwnerID),
		RoleID:    model.ProjectAdmin,
	}
	projectMemberID, err := dao.InsertOrUpdateProjectMember(projectMember)
	if err != nil {
		return false, errors.New("failed to create project member")
	}
	if projectID == 0 || projectMemberID == 0 {
		return false, errors.New("failed to create projectID memberID")
	}
	return true, nil
}

func GetProject(project model.Project, selectedFields ...string) (*model.Project, error) {
	p, err := dao.GetProject(project, selectedFields...)
	if err != nil {
		return nil, err
	}
	setDeletable(p)
	return p, nil
}

func GetProjectByName(name string) (*model.Project, error) {
	return GetProject(model.Project{Name: name, Deleted: 0}, "name", "deleted")
}

func GetProjectByID(id int64) (*model.Project, error) {
	return GetProject(model.Project{ID: id, Deleted: 0}, "id", "deleted")
}

func ProjectExists(projectName string) (bool, error) {
	query := model.Project{Name: projectName}
	project, err := dao.GetProject(query, "name")
	if err != nil {
		return false, err
	}
	return (project != nil && project.ID != 0), nil
}

func ProjectExistsByID(projectID int64) (bool, error) {
	query := model.Project{ID: projectID, Deleted: 0}
	project, err := dao.GetProject(query, "id", "deleted")
	if err != nil {
		return false, err
	}
	return (project != nil && project.Name != ""), nil
}

func UpdateProject(project model.Project, fieldNames ...string) (bool, error) {
	if project.ID == 0 {
		return false, errors.New("no Project ID provided")
	}
	_, err := dao.UpdateProject(project, fieldNames...)
	if err != nil {
		return false, err
	}
	return true, nil
}

func ToggleProjectPublic(projectID int64, public int) (bool, error) {
	return UpdateProject(model.Project{ID: projectID, Public: public}, "public")
}

func GetProjectsByUser(query model.Project, userID int64) ([]*model.Project, error) {
	projects, err := dao.GetProjectsByUser(query, userID)
	if err != nil {
		return nil, err
	}
	for i := range projects {
		setDeletable(projects[i])
	}
	return projects, nil
}

func GetPaginatedProjectsByUser(query model.Project, userID int64, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedProjects, error) {
	paged, err := dao.GetPaginatedProjectsByUser(query, userID, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		return nil, err
	}
	if paged != nil {
		for i := range paged.ProjectList {
			setDeletable(paged.ProjectList[i])
		}
	}
	return paged, nil
}

func GetProjectsByMember(query model.Project, userID int64) ([]*model.Project, error) {
	projects, err := dao.GetProjectsByMember(query, userID)
	if err != nil {
		return nil, err
	}
	for i := range projects {
		setDeletable(projects[i])
	}
	return projects, nil
}

func DeleteProject(userID, projectID int64) (bool, error) {
	project, err := GetProjectByID(projectID)
	if err != nil {
		logs.Error("Failed to delete project with ID: %d, error: %+v", projectID, err)
		return false, err
	}
	if !project.Deletable {
		logs.Error("Project %s is a builtin project that cannot be deleted.", project.Name)
		return false, utils.ErrUnprocessableEntity
	}
	members, err := GetProjectMembers(project.ID)
	if err != nil {
		return false, err
	}
	if len(members) > 1 {
		logs.Error("Project %s has member that cannot be deleted.", project.Name)
		return false, utils.ErrUnprocessableEntity
	}
	serviceList, err := dao.GetServiceData(model.ServiceStatus{ProjectID: projectID}, userID)
	if err != nil {
		logs.Error("Failed to get service data with user ID: %d, error: %+v", userID, err)
		return false, utils.ErrUnprocessableEntity
	}
	if len(serviceList) > 0 {
		logs.Error("Project %s has service deployment.", project.Name)
		return false, utils.ErrUnprocessableEntity
	}
	user, err := GetUserByID(userID)
	if err != nil {
		logs.Error("Failed to delete user with ID: %d, error: %+v", projectID, err)
		return false, err
	}
	//Delete repo in Gogits
	repoName, err := ResolveRepoName(project.Name, user.Username)
	if err != nil {
		logs.Error("Failed to resolve repo name with project name: %s, username: %s, error: %+v", project.Name, user.Username, err)
		return false, err
	}
	err = CurrentDevOps().DeleteRepo(user.Username, repoName)
	if err != nil {
		logs.Error("Failed to delete repo with repo name: %s, error: %+v", repoName, err)
		if err == utils.ErrUnprocessableEntity {
			return false, err
		}
	}
	repoPath := ResolveRepoPath(repoName, user.Username)
	err = os.RemoveAll(repoPath)
	if err != nil {
		logs.Error("Failed to remove repo path: %s, error: %+v", repoPath, err)
		if err == utils.ErrUnprocessableEntity {
			return false, err
		}
	}
	//Delete namespace in cluster
	_, err = DeleteNamespace(project.Name)
	if err != nil {
		logs.Error("Failed to delete namespace with project name: %s, error: %+v", project.Name, err)
		return false, err
	}
	project.Name = "%" + project.Name + "%"
	project.Deleted = 1
	_, err = dao.UpdateProject(*project, "name", "deleted")
	if err != nil {
		return false, err
	}
	return true, nil
}

func NamespaceExists(projectName string) (bool, error) {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Namespace()

	namespaceList, err := n.List()
	if err != nil {
		logs.Error("Failed to check namespace list in cluster", projectName)
		return false, err
	}

	for _, namespace := range (*namespaceList).Items {
		if projectName == namespace.Name {
			logs.Info("Namespace existing %+v", namespace)
			return true, nil
		}
	}
	return false, nil
}

func CreateNamespace(project *model.Project) (bool, error) {
	projectExists, err := NamespaceExists(project.Name)
	if err != nil {
		return false, err
	}
	if projectExists {
		logs.Info("Project %s already exists in cluster.", project.Name)
		return true, nil
	}

	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Namespace()

	var namespace model.Namespace
	namespace.ObjectMeta.Name = project.Name
	if project.IstioSupport {
		namespace.Labels = map[string]string{istioLabel: "enabled"}
	}
	_, err = n.Create(&namespace)
	if err != nil {
		logs.Error("Failed to create namespace: %s, error: %+v", project.Name, err)
		return false, err
	}
	logs.Info(namespace)
	return true, nil
}

func SyncNamespaceByOwnerID(userID int64) error {
	query := model.Project{OwnerID: int(userID)}
	projects, err := GetProjectsByUser(query, userID)
	if err != nil {
		return fmt.Errorf("Failed to get default projects with user ID: %d, error: %+v", userID, err)
	}

	for _, project := range projects {
		if !utils.ValidateWithPattern("project", project.Name) {
			continue
		}
		_, err = CreateNamespace(project)
		if err != nil {
			return fmt.Errorf("Failed to create namespace: %s, error: %+v", project.Name, err)
		}
	}
	return nil
}

func SyncProjectsWithK8s() error {
	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Namespace()

	namespaceList, err := n.List()
	if err != nil {
		logs.Error("Failed to check namespace list in cluster: %+v", err)
		return err
	}

	for _, namespace := range (*namespaceList).Items {
		// Skip kubernetes system namespace
		if namespace.Name == kubeNamespace || namespace.Name == istioNamespace {
			logs.Debug("Skip %s namespace", namespace.Name)
			continue
		}
		existing, err := ProjectExists(namespace.Name)
		if err != nil {
			logs.Error("Failed to check prject existing name: %s, error: %+v", namespace.Name, err)
			continue
		}
		if existing {
			logs.Info("Project existing %s", namespace.Name)
		} else {
			//Add it to projects
			var reqProject model.Project
			reqProject.Name = namespace.Name
			reqProject.OwnerID = adminUserID
			reqProject.OwnerName = adminUserName
			reqProject.Public = projectPrivate
			reqProject.CreationTime = namespace.CreationTimestamp
			reqProject.UpdateTime = namespace.CreationTimestamp
			if namespace.Labels != nil && namespace.Labels[istioLabel] == "enabled" {
				reqProject.IstioSupport = true
			}
			isSuccess, err := CreateProject(reqProject)
			if err != nil {
				logs.Error("Failed to create project name: %s, error: %+v", reqProject.Name, err)
				// Still can work
				continue
			}
			if !isSuccess {
				logs.Error("Failed to create project name: %s, error: %+v", reqProject.Name, err)
				// Still can work
				continue
			}
			err = CurrentDevOps().CreateRepoAndJob(adminUserID, reqProject.Name)
			if err != nil {
				logs.Error("Failed create repo and job with project name: %s, error: %+v", reqProject.Name, err)
			}
		}
		// Sync the helm release on this project namespace
		err = SyncHelmReleaseWithK8s(namespace.Name)
		if err != nil {
			logs.Error("Failed to sync helm service with project name: %s, error: %+v", namespace.Name, err)
			// Still can work
		}

		// Sync the services in this project namespace
		err = SyncServiceWithK8s(namespace.Name)
		if err != nil {
			logs.Error("Failed to sync service with project name: %s, error: %+v", namespace.Name, err)
			// Still can work
		}

		// Sync the autoscale hpa in this project namespace
		err = SyncAutoScaleWithK8s(namespace.Name)
		if err != nil {
			logs.Error("Failed to sync autoscale rule with project name: %s, error: %+v", namespace.Name, err)
			// Still can work
		}

		// Sync the services in this project namespace
		err = SyncJobWithK8s(namespace.Name)
		if err != nil {
			logs.Error("Failed to sync job with project name: %s, error: %+v", namespace.Name, err)
			// Still can work
		}
	}
	return err
}

func DeleteNamespace(nameSpace string) (bool, error) {
	namespaceExists, err := NamespaceExists(nameSpace)
	if err != nil {
		return false, err
	}
	if !namespaceExists {
		logs.Info("Namespace %s not exists in cluster.", nameSpace)
		return true, nil
	}

	var config k8sassist.K8sAssistConfig
	config.KubeConfigPath = kubeConfigPath()
	k8sclient := k8sassist.NewK8sAssistClient(&config)
	n := k8sclient.AppV1().Namespace()

	err = n.Delete(nameSpace)
	if err != nil {
		logs.Error("Failed to delete namespace %s", nameSpace)
		return false, err
	}
	return true, nil
}

func setDeletable(p *model.Project) {
	if p == nil {
		return
	}
	for _, name := range undeletableNamespaces {
		if name == p.Name {
			p.Deletable = false
			return
		}
	}
	p.Deletable = true
}
