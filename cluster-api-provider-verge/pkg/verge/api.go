package vergeapis

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha2"
	"sigs.k8s.io/cluster-api/exp/runtime/client"
)

// getOwnerMachine returns the owner machine object: TODO
func GetOwnerMachine(ctx context.Context, client client.Client, owner metav1.ObjectMeta) (*clusterv1.Machine, error) {
	return nil, nil
}

// createMachineInVergeInfra creates a machine in the provider and returns the providerID: TODO
func CreateMachineInVergeInfra(machine *clusterv1.Machine, userData []byte) (string, error) {
	return "", nil
}

// deleteMachineInVergeInfra deletes a machine in the provider: TODO
func DeleteMachineInVergeInfra(machine *clusterv1.Machine) error {
	return nil
}
