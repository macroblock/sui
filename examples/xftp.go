package main

import "io"
import "github.com/pkg/sftp"
import "fmt"

type IFtp interface {
	//Close()
	FileSize(path string) (int64, error)
	Delete(path string) error
	Rename(from, to string) error
	StorFrom(path string, r io.Reader, offset uint64) error
	Quit() error
}

type TSftp struct {
	client *sftp.Client
}

func (o *TSftp) Delete(path string) error {
	return o.client.Remove(path)
}

func (o *TSftp) Rename(from, to string) error {
	return o.client.Rename(from, to)
}

func (o *TSftp) Quit() error {
	return o.client.Close()
}

func (o *TSftp) FileSize(path string) (int64, error) {
	stat, err := o.client.Stat(path)
	if err != nil {
		return -1, err
	}
	return stat.Size(), nil
}

func (o *TSftp) StorFrom(path string, r io.Reader, offset uint64) error {
	// conn, err := c.cmdDataConnFrom(offset, "STOR %s", path)
	// if err != nil {
	// 	return err
	// }
	f, err := o.client.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	//_, err = io.Copy(f, r)
	offs, err := f.Seek(int64(offset), 0)
	if err != nil {
		return err
	}
	if offs != int64(offset) {
		return fmt.Errorf("Sftp Seek() problem (custom error). Search %v, but return %v", int64(offset), offs)
	}
	_, err = f.ReadFrom(r)
	if err != nil {
		return err
	}
	// _, _, err = c.conn.ReadResponse(StatusClosingDataConnection)
	// return err
	return nil
}
