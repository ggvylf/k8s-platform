package model

import "gorm.io/gorm"

// CREATE TABLE workflow (
//     ID INT(10) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
//     CreatedAt datetime DEFAULT NULL,
//     UpdatedAt datetime DEFAULT NULL,
//     DeletedAt datetime DEFAULT NULL,
//     Name VARCHAR(32) DEFAULT NULL,
//     Namespace VARCHAR(32) DEFAULT NULL,
//     Replicas INT DEFAULT NULL,
//     Deployment VARCHAR(32) DEFAULT NULL,
//     Service VARCHAR(32) DEFAULT NULL,
//     Ingress VARCHAR(32) DEFAULT NULL,
//     Type VARCHAR(32) DEFAULT NULL
// );

type Workflow struct {
	// 默认字段
	gorm.Model

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	// svc的资源类型 clusterip nodeport ingress
	Type string `json:"service_type",gorm:"cloumn:type"`
}

func (*Workflow) TableName() string {
	return "workflow"
}
