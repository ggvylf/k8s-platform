package service

import (
	"k8s-platform/dao"
	"k8s-platform/model"
)

type workflow struct{}

var Workflow workflow

// 整合了相关资源的参数
// type表示svc资源的类型 Ingress nodePort ClusterIP等
type WorkflowCreate struct {
	Name          string                 `json:"name"`
	Namespace     string                 `json:"namespace"`
	Replicas      int32                  `json:"replicas"`
	Image         string                 `json:"image"`
	Label         map[string]string      `json:"label"`
	Cpu           string                 `json:"cpu"`
	Memory        string                 `json:"memory"`
	ContainerPort int32                  `json:"container_port"`
	HealthCheck   bool                   `json:"health_check"`
	HealthPath    string                 `json:"health_path"`
	Type          string                 `json:"type"`
	Port          int32                  `json:"port"`
	NodePort      int32                  `json:"node_port"`
	Hosts         map[string][]*HttpPath `json:"hosts"`
}

func (w *workflow) GetList(name, namespace string, page, limit int) (data *dao.WorkflowResp, err error) {
	data, err = dao.Workflow.GetList(name, namespace, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) GetById(id int) (data *model.Workflow, err error) {
	data, err = dao.Workflow.GetById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *workflow) CreateWorkFlow(data *WorkflowCreate) (err error) {
	var ingressName string

	// 判断type
	if data.Type == "Ingress" {
		ingressName = getIngressName(data.Name)
	} else {
		ingressName = ""
	}

	// 组装workflow
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    getServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}

	// 写入db
	err = dao.Workflow.Add(workflow)
	if err != nil {
		return err
	}

	err = createWorkflowRes(data)
	if err != nil {
		return err
	}
	return err
}

// 根据workflow 创建对应的deploy svc ingress
func createWorkflowRes(data *WorkflowCreate) (err error) {
	// deploy
	dc := &DeployCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Replicas:      data.Replicas,
		Image:         data.Image,
		Label:         data.Label,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
	}

	err = Deployment.CreateDeployment(dc)
	if err != nil {
		return err
	}

	// svc
	var serviceType string
	if data.Type != "Ingress" {
		serviceType = data.Type
	} else {
		serviceType = "ClusterIP"
	}

	sc := &ServiceCreate{
		Name:          getServiceName(data.Name),
		Namespace:     data.Namespace,
		Type:          serviceType,
		ContainerPort: data.ContainerPort,
		Port:          data.Port,
		NodePort:      data.NodePort,
		Label:         data.Label,
	}
	if err := Service.CreateService(sc); err != nil {
		return err
	}

	// ingress
	var ic *IngressCreate
	if data.Type == "Ingress" {
		ic = &IngressCreate{
			Name:      getIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(ic)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *workflow) DeleteWorkflowById(id int) (err error) {

	// 从db中获取workflow信息
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}

	// 删除workflow资源
	err = delWorkflowRes(workflow)
	if err != nil {
		return err
	}

	// 从db中删除workflow
	err = dao.Workflow.DeleteByID(id)
	if err != nil {
		return err
	}
	return
}

func delWorkflowRes(workflow *model.Workflow) (err error) {
	// deploy
	err = Deployment.DeleteDeployment(workflow.Name, workflow.Namespace)
	if err != nil {
		return err
	}

	// svc
	err = Service.DeleteService(getServiceName(workflow.Name), workflow.Namespace)
	if err != nil {
		return err
	}

	// ingress
	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(getIngressName(workflow.Name), workflow.Namespace)
		if err != nil {
			return err
		}
	}

	return nil
}

// 资源名称规范，在workflow的name后边加对应资源的后缀
// svc

func getServiceName(workflowName string) (serviceName string) {
	return workflowName + "-svc"
}

// ingress
func getIngressName(workflowName string) (ingressName string) {
	return workflowName + "-ing"
}
