package windows

import (
	"context"
	"github.com/risor-io/risor/arg"
	"github.com/risor-io/risor/object"
	"github.com/risor-io/risor/op"
	"golang.org/x/sys/windows"
	"strings"
)

func GetLogicalDriveStrings(ctx context.Context, args ...object.Object) object.Object {
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		return object.NewError(err)
	}
	drives := make([]string, 0, 'Z'-'A')
	for i := 'A'; i <= 'Z'; i++ {
		ok := (mask & (0x01 << (i - 'A'))) != 0
		if !ok {
			continue
		}
		drives = append(drives, string(i)+":\\")
	}
	return object.NewStringList(drives)
}

type DriveType uint32

func (d DriveType) valueToName() string {
	switch d {
	case 2:
		return "removable"
	case 3:
		return "fixed"
	case 4:
		return "remote"
	case 5:
		return "cdrom"
	case 6:
		return "ramdisk"
	}
	return ""
}

func (d DriveType) Inspect() string {
	return d.valueToName()
}

func (d DriveType) Interface() interface{} {
	return d
}

func (d DriveType) Equals(other object.Object) object.Object {
	return object.NewBool(d == other)
}

func (d DriveType) GetAttr(name string) (object.Object, bool) {
	switch name {
	case "is_removable", "is_fixed", "is_remote", "is_cdrom", "is_ramdisk":
		builtinName := "windows" + "." + string(d.Type()) + "." + name
		return object.NewBuiltin(builtinName, func(ctx context.Context, args ...object.Object) object.Object {
			return object.NewBool(strings.TrimPrefix(name, "is_") == d.valueToName())
		}), true
	default:
		return object.TypeErrorf("type error: windows.%v object has no attribute %q", d.Type(), name), false
	}
}

func (d DriveType) SetAttr(name string, value object.Object) error {
	return object.TypeErrorf("type error: windows.%v object has no attribute %q", d.Type(), name)
}

func (d DriveType) IsTruthy() bool {
	return true
}

func (d DriveType) RunOperation(opType op.BinaryOpType, right object.Object) object.Object {
	return object.TypeErrorf("type error: unsupported operation for windows.%v: %v", d.Type(), opType)
}

func (d DriveType) Cost() int {
	return 0
}

func (d DriveType) Type() object.Type {
	return "drive_type"
}

func GetDriveType(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("windows.get_drive_type", 1, args); err != nil {
		return err
	}
	path, rerr := object.AsString(args[0])
	if rerr != nil {
		return rerr
	}
	pathName, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return object.NewError(err)
	}
	driveType := windows.GetDriveType(pathName)
	if driveType == 0 {
		return object.Errorf("%q: cannot determine drive type", pathName)
	} else if driveType == 1 {
		return object.Errorf("%q: invalid root path", pathName)
	}
	return DriveType(driveType)
}
