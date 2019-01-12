package scp

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/Zeng1999/motion/conn/ssh"
)

type ScpTrans struct {
	From string `json:"from"`
	Des  string `json:"des"`

	conn *ssh.SSHConn
}

func NewScpTrans(con *ssh.SSHConn) *ScpTrans {
	return &ScpTrans{
		conn: con,
	}
}

func (ts *ScpTrans) Init() error {
	if ts.conn.GetClient() != nil {
		return nil
	}
	return ts.conn.Dial()
}

func (ts *ScpTrans) Run() error {
	var err error

	err = ts.Init()
	if err != nil {
		return err
	}

	ses, err := ts.conn.NewSession()
	if err != nil {
		return err
	}
	defer ses.Close()

	data, err := ioutil.ReadFile(ts.From)
	if err != nil {
		return err
	}

	w, err := ses.StdinPipe()
	if err != nil {
		return err
	}

	go routine(w, data, ts.Des)

	return ses.Run("/usr/bin/scp -qrt ./")
}

func routine(w io.WriteCloser, data []byte, des string) {
	defer w.Close()
	fmt.Fprintln(w, "C0644", len(data), des)
	w.Write(data)
	fmt.Fprint(w, "\x00")
}

func (ts *ScpTrans) SetFrom(from string) {
	ts.From = from
}

func (ts *ScpTrans) SetDes(destination string) {
	ts.Des = destination
}

func (ts *ScpTrans) Close() error {
	if ts.conn == nil {
		return nil
	}
	return ts.conn.Close()
}
