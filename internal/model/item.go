package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Item struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	ID       string `json:"id"`
	Name     string `json:"name"`
	Revision int    `json:"revision"`
	Status   struct {
		LastReconciledRevision int `json:"lastReconciledRevision"`
	} `json:"status"`
}

func (i Item) IsPendingChanges() bool {
	return i.Revision != i.Status.LastReconciledRevision
}

// func init() {
// 	SchemeBuilder.Register(&Item{})
// }
