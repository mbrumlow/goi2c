package i2c

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type I2C struct {
	dev *os.File
}

const (
	i2cSlave = 0x0703
	i2cSMBus = 0x0720

	i2cSMBusRead  = uint8(1)
	i2cSMBusWrite = uint8(0)

	i2cSMBusByteData = uint32(2)
)

type i2c_smbus_ioctl_data_byte struct {
	read_write uint8
	command    uint8
	size       uint32
	data       *uint8
}

func New(addr uint8, bus int) (*I2C, error) {

	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	if err := ioctl(f.Fd(), i2cSlave, uintptr(addr)); err != nil {
		return nil, err
	}

	return &I2C{dev: f}, nil
}

func (i2c *I2C) ReadRegister(r uint8) (uint8, error) {

	v := uint8(0)

	args := i2c_smbus_ioctl_data_byte{
		read_write: i2cSMBusRead,
		command:    r,
		size:       i2cSMBusByteData,
		data:       &v,
	}

	if err := ioctl(i2c.dev.Fd(), i2cSMBus, uintptr(unsafe.Pointer(&args))); err != nil {
		return v, err
	}

	return v, nil

}

func (i2c *I2C) WriteRegister(r, v uint8) error {

	args := i2c_smbus_ioctl_data_byte{
		read_write: i2cSMBusWrite,
		command:    r,
		size:       i2cSMBusByteData,
		data:       &v,
	}

	if err := ioctl(i2c.dev.Fd(), i2cSMBus, uintptr(unsafe.Pointer(&args))); err != nil {
		return err
	}

	return nil
}

func (i2c *I2C) Close() error {
	return i2c.dev.Close()
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, cmd, arg)
	if err != 0 {
		return err
	}
	return nil
}
