//go:build windows

package utils

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type JobObject struct {
	handle windows.Handle
}

// Link process with the main process so it is killed when the main process is killed
func NewJobObject() (*JobObject, error) {
	h, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return nil, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	_, err = windows.SetInformationJobObject(
		h,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	if err != nil {
		windows.CloseHandle(h)
		return nil, err
	}

	return &JobObject{handle: h}, nil
}

func (j *JobObject) AddProcess(p *os.Process) error {
	// Open a real Windows handle for the process
	h, err := windows.OpenProcess(windows.PROCESS_ALL_ACCESS, false, uint32(p.Pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(h)

	return windows.AssignProcessToJobObject(j.handle, h)
}

func (j *JobObject) Close() error {
	return windows.CloseHandle(j.handle)
}
