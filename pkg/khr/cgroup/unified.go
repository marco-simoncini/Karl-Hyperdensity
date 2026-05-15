package cgroup

// UnifiedCgroupMount is the default unified hierarchy mount on Linux.
const UnifiedCgroupMount = "/sys/fs/cgroup"

// DefaultScannedRoot returns the default read-only scan root for cgroup discovery.
func DefaultScannedRoot() string {
	return UnifiedCgroupMount
}
