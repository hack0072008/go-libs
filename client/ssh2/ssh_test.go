/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2020 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package ssh2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var ss *SSH

func init() {
	port, err := strconv.Atoi(os.Getenv("SSH_PORT"))
	utilruntime.Must(err)
	ss, _ = New(&Config{
		Host:     os.Getenv("SSH_HOST"),
		Port:     port,
		User:     os.Getenv("SSH_USER"),
		Password: os.Getenv("SSH_PASSWORD"),
	})
}

func TestCombinedOutput(t *testing.T) {
	cmd1 := fmt.Sprintf(`sh /root/exp3.sh`)

	t.Logf("cmd1: %+v", cmd1)
	output, err := ss.CombinedOutput(cmd1)
	assert.Nil(t, err)
	t.Logf("output1: %+v", string(output))
}

func TestSudo(t *testing.T) {
	output, err := ss.CombinedOutput("type expect")
	// output, err := s.CombinedOutput("cat /etc/fstab|grep UUID|awk '{print $1}'|awk -F'=' '{print $2}'")
	// output, err := s.CombinedOutput("type scp")
	// output, err := s.CombinedOutput("type tar")
	// output, err := s.CombinedOutput("type sh")
	// output, err := s.CombinedOutput("parted /dev/vdb mkpart primary 0 100%")
	assert.Nil(t, err)
	// assert.Equal(t, "root", strings.TrimSpace(string(output)))
	t.Logf("ssh sudo host: %+v port: %+v user: %+v succ, output: %+v", ss.Host, ss.Port, ss.User, string(output))
}

func TestQuote(t *testing.T) {
	output, err := ss.CombinedOutput(`echo "a" 'b'`)
	assert.Nil(t, err)
	assert.Equal(t, "a b", strings.TrimSpace(string(output)))
	t.Logf("ssh quote host: %+v port: %+v user: %+v succ, output: %+v", ss.Host, ss.Port, ss.User, string(output))
}

func TestWriteFile(t *testing.T) {
	// data := []byte("Hello")
	data := []byte(fmt.Sprintf(`
#!/bin/bash
expect <<EOF
   set timeout %d
   spawn scp %s %s@%s:%s
   expect {
       "yes/no" { send "yes\r";exp_continue }
       "password" { send "%s\r" }
   }
   expect "]#" { send "exit\r" } expect eof
EOF
`, 60, "/root/ahc_x86", "root", "123.249.87.110", "/cpaas/", "Pass@Ecs"))
	dst := "/root/exp3.sh"

	err := ss.WriteFile(bytes.NewBuffer(data), dst)
	assert.Nil(t, err)

	output, err := ss.ReadFile(dst)
	assert.Nil(t, err)
	assert.Equal(t, data, output)
}

func TestCoppyFile(t *testing.T) {
	src := os.Args[0]
	srcData, err := ioutil.ReadFile(src)
	assert.Nil(t, err)

	dst := "/tmp/test"
	err = ss.CopyFile(src, dst)
	assert.Nil(t, err)

	output, err := ss.ReadFile(dst)
	assert.Nil(t, err)

	assert.Equal(t, srcData, output)
}

func TestExist(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			"exist",
			args{
				filename: "/tmp",
			},
			true,
			false,
		},
		{
			"not exist",
			args{
				filename: "/tmpfda",
			},
			false,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ss.Exist(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exist() got = %v, want %v", got, tt.want)
			}
		})
	}
}
