package checkers

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/efureev/health-checker"
)

func CheckingHDDFn(path string, size uint64) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		disk := DiskUsage(path)

		result.Info.Text = disk.String()

		if !disk.HasFree(size) {
			result.AddError(errors.New(`lack of disk space`))
			return
		}

		return
	}
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All  uint64
	Used uint64
	Free uint64
}

func (d DiskStatus) AllString() string {
	return fmt.Sprintf("All: %.2f GB", float64(d.All)/float64(GB))
}

func (d DiskStatus) UsedString() string {
	return fmt.Sprintf("Used: %.2f GB", float64(d.Used)/float64(GB))
}

func (d DiskStatus) FreeString() string {
	return fmt.Sprintf("Free: %.2f GB", float64(d.Free)/float64(GB))
}

func (d DiskStatus) HasFree(size uint64) bool {
	return d.Free > size
}

func (d DiskStatus) String() string {
	return fmt.Sprintf("%s\n%s\n%s\n", d.AllString(), d.UsedString(), d.FreeString())
}

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
