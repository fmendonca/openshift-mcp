package rbac

import (
	"fmt"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
)

func formatRolesList(roles *rbacv1.RoleList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total Roles: %d\n\n", len(roles.Items)))

	for _, role := range roles.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", role.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", role.Namespace))
		sb.WriteString(fmt.Sprintf("Rules: %d\n", len(role.Rules)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatRoleBindingsList(roleBindings *rbacv1.RoleBindingList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total RoleBindings: %d\n\n", len(roleBindings.Items)))

	for _, rb := range roleBindings.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", rb.Name))
		sb.WriteString(fmt.Sprintf("Namespace: %s\n", rb.Namespace))
		sb.WriteString(fmt.Sprintf("Role: %s\n", rb.RoleRef.Name))
		sb.WriteString(fmt.Sprintf("Subjects: %d\n", len(rb.Subjects)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatClusterRolesList(clusterRoles *rbacv1.ClusterRoleList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total ClusterRoles: %d\n\n", len(clusterRoles.Items)))

	for _, cr := range clusterRoles.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", cr.Name))
		sb.WriteString(fmt.Sprintf("Rules: %d\n", len(cr.Rules)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}

func formatClusterRoleBindingsList(clusterRoleBindings *rbacv1.ClusterRoleBindingList) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Total ClusterRoleBindings: %d\n\n", len(clusterRoleBindings.Items)))

	for _, crb := range clusterRoleBindings.Items {
		sb.WriteString(fmt.Sprintf("Name: %s\n", crb.Name))
		sb.WriteString(fmt.Sprintf("Role: %s\n", crb.RoleRef.Name))
		sb.WriteString(fmt.Sprintf("Subjects: %d\n", len(crb.Subjects)))
		sb.WriteString("\n---\n\n")
	}

	return sb.String()
}
