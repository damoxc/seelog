// Copyright (c) 2012 - Cloud Instruments Co. Ltd.
// 
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"os"
	"errors"
	"runtime"
	"strings"
	"path/filepath"
	"time"
)

// Represents runtime caller context
type logContext struct {
	funcName  string
	shortPath string
	fullPath  string
	fileName string
	callTime time.Time
}

// Returns context of the caller
func currentContext() (*logContext, error) {
	return specificContext(1)
}

// Returns context of the function with placed "skip" stack frames of the caller
// If skip == 0 then behaves like currentContext
func specificContext(skip int) (*logContext, error) {
	if skip < 0 {
		return nil, errors.New("Can not skip negative stack frames")
	}
	
	fullPath, shortPath, function, err := extractCallerInfo(skip + 2)
	if err != nil {
		return nil, err
	}
	_, fileName := filepath.Split(fullPath)
	return &logContext{function, shortPath, fullPath, fileName, time.Now()}, nil
}

func (context *logContext) Func() string {
	return context.funcName
}

func (context *logContext) ShortPath() string {
	return context.shortPath
}

func (context *logContext) FullPath() string {
	return context.fullPath
}

func (context *logContext) FileName() string {
	return context.fileName
}

func (context *logContext) CallTime() time.Time {
	return context.callTime;
}

var workingDir = ""

func init() {
	setWorkDir()
}

func setWorkDir() {
	workDir, workingDirError := os.Getwd()
	if workingDirError != nil {
		workingDir = string(os.PathSeparator)
		return
	}

	workingDir = workDir + string(os.PathSeparator)
}

func extractCallerInfo(skip int) (fullPath string, shortPath string, funcName string,err error) {
	pc, fullPath, _, ok := runtime.Caller(skip)

	if !ok {
		return "", "", "", errors.New("Error during runtime.Caller")
	}

	//TODO:Currently fixes bug in weekly.2012-03-13+: Caller returns incorrect separators
	//Delete later
	
	fullPath = strings.Replace(fullPath, "\\", string(os.PathSeparator), -1)
	fullPath = strings.Replace(fullPath, "/", string(os.PathSeparator), -1)

	
	if strings.HasPrefix(fullPath, workingDir) {
		shortPath = fullPath[len(workingDir):]
	} else {
		shortPath = fullPath
	}

	funName := runtime.FuncForPC(pc).Name()
	var functionName string
	if strings.HasPrefix(funName, workingDir) {
		functionName = funName[len(workingDir):len(funName)]
	} else {
		functionName = funName
	}

	return fullPath, shortPath, functionName, nil
}
