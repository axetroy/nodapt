package node

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

type CachedNode struct {
	Version  string
	FilePath string
}

// 按照版本号升序排序
type ByVersion []CachedNode

// 实现 sort.Interface 的 Len() 方法
func (a ByVersion) Len() int {
	return len(a)
}

// 实现 sort.Interface 的 Less() 方法，定义排序逻辑
func (a ByVersion) Less(i, j int) bool {
	ver1, err1 := semver.NewVersion(a[i].Version)
	ver2, err2 := semver.NewVersion(a[j].Version)

	// In case of parsing errors, fall back to string comparison
	if err1 != nil || err2 != nil {
		return a[i].Version < a[j].Version
	}

	return ver1.LessThan(ver2) // 根据 Version 字段升序排序
}

// 实现 sort.Interface 的 Swap() 方法，交换元素
func (a ByVersion) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func GetCachedVersions(nodaptDir string) ([]CachedNode, error) {
	list := make([]CachedNode, 0)

	nodeDir := filepath.Join(nodaptDir, "node")

	if stat, err := os.Stat(nodeDir); err != nil {
		if os.IsNotExist(err) {
			return list, nil
		}

		return nil, errors.WithStack(err)
	} else if !stat.IsDir() {
		return nil, errors.Errorf("node directory is not a directory")
	}

	entries, err := os.ReadDir(nodeDir)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, file := range entries {
		fName := file.Name()

		if file.IsDir() && strings.HasPrefix(fName, "node-v") {
			n := strings.SplitN(fName, "-", -1)

			version := n[1]

			list = append(list, CachedNode{
				Version:  version,
				FilePath: filepath.Join(nodeDir, fName),
			})
		}
	}

	// Sort versions in ascending order
	sort.Sort(ByVersion(list))

	return list, nil
}
