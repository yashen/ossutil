package lib 

import (
    "fmt"
    "strconv"
    "os"
    "time"
    "strings"

    . "gopkg.in/check.v1"
)

func (s *OssutilCommandSuite) TestCPObject(c *C) {
    bucket := bucketNameCP 
    destBucket := bucketNameNotExist 

    // put object
    s.createFile(uploadFileName, content, c)
    object := "TestCPObject_cp" 
    s.putObject(bucket, object, uploadFileName, c)

    // get object
    s.getObject(bucket, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, content)

    // modify uploadFile content
    data := "欢迎使用ossutil"
    s.createFile(uploadFileName, data, c)

    time.Sleep(sleepTime)
    // put to exist object
    s.putObject(bucket, object, uploadFileName, c)

    // get to exist file
    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // get without specify dest file 
    s.getObject(bucket, object, ".", c)
    str = s.readFile(object, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(object)

    // put without specify dest object 
    data1 := "put without specify dest object"
    s.createFile(uploadFileName, data1, c)
    s.putObject(bucket, "", uploadFileName, c)
    s.getObject(bucket, uploadFileName, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data1)

    // get to file in not exist directory
    notexistdir := "NOTEXISTDIR"
    s.getObject(bucket, object, notexistdir + string(os.PathSeparator) + downloadFileName, c)
    str = s.readFile(notexistdir + string(os.PathSeparator) + downloadFileName, c) 
    c.Assert(str, Equals, data)
    _ = os.RemoveAll(notexistdir)

    // copy file
    destObject := "TestCPObject_destObject"
    s.copyObject(bucket, object, bucket, destObject, c)

    objectStat := s.getStat(bucket, destObject, c)
    c.Assert(objectStat[StatACL], Equals, "default")
    
    // get dest file
    filePath := downloadFileName + "1" 
    s.getObject(bucket, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // put to not exist bucket
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(destBucket, object), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(destBucket, object), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // get not exist object
    showElapse, err = s.rawCP(CloudURLToString(bucket, "notexistobject"), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy to not exist bucket
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, destObject), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // corse bucket copy
    destBucket = bucketNameDest

    s.copyObject(bucket, object, destBucket, destObject, c)

    s.getObject(destBucket, destObject, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data)
    _ = os.Remove(filePath)

    // copy single object in directory, test the name of dest object 
    srcObject := "a/b/c/d/e"
    s.putObject(bucket, srcObject, uploadFileName, c)
    time.Sleep(time.Second)

    s.copyObject(bucket, srcObject, destBucket, "", c)

    s.getObject(destBucket, "e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucket, srcObject, destBucket, "a/", c)

    s.getObject(destBucket, "a/e", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    s.copyObject(bucket, srcObject, destBucket, "a", c)

    s.getObject(destBucket, "a", filePath, c)
    str = s.readFile(filePath, c)
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)

    // copy without specify dest object
    s.copyObject(bucket, uploadFileName, destBucket, "", c)
    s.getObject(destBucket, uploadFileName, filePath, c)
    str = s.readFile(filePath, c) 
    c.Assert(str, Equals, data1)
    _ = os.Remove(filePath)
}

func (s *OssutilCommandSuite) TestErrorCP(c *C) {
    bucket := bucketNameExist 

    // error src_url
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString("", ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString("", bucket), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", bucket), downloadFileName, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString("", ""), downloadFileName, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, "a", true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // miss argc
    showElapse, err = s.rawCP(CloudURLToString("", bucket), "", true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy self
    object := "testobject"
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, object), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(bucket, object), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err checkpoint dir, conflict with config file
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, object), false, true, true, DefaultBigFileThreshold, configFile)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestUploadErrSrc(c *C) {
    srcBucket := bucketNamePrefix + "uploadsrc"
    destBucket := bucketNameNotExist 
    command := "cp"
    args := []string{uploadFileName, CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "")}
    str := ""
    ok := true
    cpDir := CheckpointDir
    thre := strconv.FormatInt(DefaultBigFileThreshold, 10)
    routines := strconv.Itoa(Routines)
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &configFile,
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestBatchCPObject(c *C) {
    bucket := bucketNameBCP

    // create local dir
    dir := "TestBatchCPObject"
    err := os.MkdirAll(dir, 0755)
    c.Assert(err, IsNil)

    // upload empty dir miss recursive
    showElapse, err := s.rawCP(dir, CloudURLToString(bucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // upload empty dir
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)

    // head object 
    showElapse, err = s.rawGetStat(bucket, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawGetStat(bucket, dir + "/")
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)

    // create dir in dir 
    dir = "TestBatchCPObject_dir"
    subdir := "SUBDIR"
    err = os.MkdirAll(dir + string(os.PathSeparator) + subdir, 0755)
    c.Assert(err, IsNil)

    // upload dir    
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true) 

    // remove object
    s.removeObjects(bucket, subdir + "/", false, true, c)

    // create file in dir
    num := 3 
    filePaths := []string{subdir + "/"}
    for i := 0; i < num; i++ {
        filePath := fmt.Sprintf("TestBatchCPObject_%d", i) 
        s.createFile(dir + "/" + filePath, fmt.Sprintf("测试文件：%d内容", i), c)
        filePaths = append(filePaths, filePath)
    }

    // upload files
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    
    time.Sleep(7*time.Second)

    // get files
    downDir := "下载目录"
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    for _, filePath := range filePaths {
        _, err := os.Stat(downDir + "/" + filePath)
        c.Assert(err, IsNil)
    }

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    // get to exist files
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, false, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downDir, true, false, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _, err = os.Stat(downDir)
    c.Assert(err, IsNil)
    //c.Assert(f.ModTime(), Equals, f1.ModTime())

    // copy files
    destBucket := bucketNameNotExist 
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, "123"), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    destBucket = bucketNameDest

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, "123"), true, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(7*time.Second)

    for _, filePath := range filePaths {
        s.getStat(destBucket, "123" + filePath, c)
    }

    // remove dir
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(downDir)
}

func (s *OssutilCommandSuite) TestCPObjectUpdate(c *C) {
    bucket := bucketNameExist 
    s.removeObjects(bucket, "", true, true, c)
    time.Sleep(2*7*time.Second) 

    // create older file and newer file
    oldData := "old data"
    oldFile := "oldFile"
    newData := "new data"
    newFile := "newFile"
    s.createFile(oldFile, oldData, c)
    time.Sleep(7*time.Second)
    s.createFile(newFile, newData, c)

    // put newer object
    object := "testobject"
    s.putObject(bucket, object, newFile, c)

    // get object
    s.getObject(bucket, object, downloadFileName, c)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // put old object with update
    showElapse, err := s.rawCP(oldFile, CloudURLToString(bucket, object), false, false, true, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    time.Sleep(7*time.Second)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucket, object), false, true, true, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    showElapse, err = s.rawCP(oldFile, CloudURLToString(bucket, object), false, false, false, DefaultBigFileThreshold, CheckpointDir)  
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(bucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, newData)

    // get object with update 
    // modify downloadFile
    time.Sleep(1)
    downData := "download file has been modified locally"
    s.createFile(downloadFileName, downData, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, false, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, true, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, downData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, false, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    str = s.readFile(downloadFileName, c)
    c.Assert(str, Equals, downData)

    // copy object with update
    destBucket := bucketNameDest 

    destData := "data for dest bucket"
    destFile := "destFile"
    s.createFile(destFile, destData, c)
    s.putObject(destBucket, object, destFile, c) 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, false, true, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, true, true, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, object, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, destData)

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, object), false, false, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), CloudURLToString(destBucket, ""), true, false, false, DefaultBigFileThreshold, CheckpointDir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(oldFile)
    _ = os.Remove(newFile)
    _ = os.Remove(destFile)
}

func (s *OssutilCommandSuite) TestResumeCPObject(c *C) { 
    var threshold int64
    threshold = 1
    cpDir := "checkpoint目录" 

    bucket := bucketNameExist 

    data := "resume cp"
    s.createFile(uploadFileName, data, c)

    // put object
    object := "object" 
    showElapse, err := s.rawCP(uploadFileName, CloudURLToString(bucket, object), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // get object
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str := s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    s.createFile(downloadFileName, "-------long file which must be truncated by cp file------", c)
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), downloadFileName, false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)

    // copy object
    destBucket := bucketNameDest 
    s.putBucket(destBucket, c)

    destObject := "destObject" 

    showElapse, err = s.rawCP(CloudURLToString(bucket, object), CloudURLToString(destBucket, destObject), false, true, false, threshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    s.getObject(destBucket, destObject, downloadFileName, c)
    str = s.readFile(downloadFileName, c) 
    c.Assert(str, Equals, data)
}

func (s *OssutilCommandSuite) TestCPMulitSrc(c *C) {
    bucket := bucketNameExist 

    // upload multi file 
    file1 := uploadFileName + "1"
    s.createFile(file1, file1, c)
    file2 := uploadFileName + "2"
    s.createFile(file2, file2, c)
    showElapse, err := s.rawCPWithArgs([]string{file1, file2, CloudURLToString(bucket, "")}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _ = os.Remove(file1)
    _ = os.Remove(file2)

    // download multi objects
    object1 := "object1"
    s.putObject(bucket, object1, uploadFileName, c)
    object2 := "object2"
    s.putObject(bucket, object2, uploadFileName, c)
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucket, object1), CloudURLToString(bucket, object2), "../"}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // copy multi objects
    destBucket := bucketNameDest 
    showElapse, err = s.rawCPWithArgs([]string{CloudURLToString(bucket, object1), CloudURLToString(bucket, object2), CloudURLToString(destBucket, "")}, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrUpload(c *C) {
    // src file not exist
    bucket := bucketNameExist 
    
    showElapse, err := s.rawCP("notexistfile", CloudURLToString(bucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // create local dir
    dir := "上传目录"
    err = os.MkdirAll(dir, 0755)
    c.Assert(err, IsNil)
    cpDir := dir + string(os.PathSeparator) + CheckpointDir 
    showElapse, err = s.rawCP(dir, CloudURLToString(bucket, ""), true, true, true, DefaultBigFileThreshold, cpDir)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    // err object name
    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, "/object"), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(uploadFileName, CloudURLToString(bucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    subdir := dir + string(os.PathSeparator) + "subdir"
    err = os.MkdirAll(subdir, 0755)
    c.Assert(err, IsNil)

    showElapse, err = s.rawCP(subdir, CloudURLToString(bucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(subdir)
}

func (s *OssutilCommandSuite) TestErrDownload(c *C) {
    bucket := bucketNameExist 
 
    object := "object"
    s.putObject(bucket, object, uploadFileName, c)

    // download to dir, but dir exist as a file
    showElapse, err := s.rawCP(CloudURLToString(bucket, object), configFile + string(os.PathSeparator), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // batch download without -r
    showElapse, err = s.rawCP(CloudURLToString(bucket, ""), downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // download to file in not exist dir
    showElapse, err = s.rawCP(CloudURLToString(bucket, object), configFile + string(os.PathSeparator) + downloadFileName, false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestErrCopy(c *C) {
    srcBucket := bucketNameExist 

    destBucket := bucketNameDest 

    // batch copy without -r
    showElapse, err := s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // error src object name
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, "/object"), CloudURLToString(destBucket, ""), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    // err dest object
    object := "object"
    s.putObject(srcBucket, object, uploadFileName, c)
    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, DefaultBigFileThreshold, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, "/object"), false, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    showElapse, err = s.rawCP(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, "/object"), true, true, false, 1, CheckpointDir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestPreparePartOption(c *C) {
    partSize, routines := copyCommand.preparePartOption(100000000000)
    c.Assert(partSize, Equals, int64(250000000))
    c.Assert(routines, Equals, 5)

    partSize, routines = copyCommand.preparePartOption(80485760)
    c.Assert(partSize, Equals, int64(12816225))
    c.Assert(routines, Equals, 2)

    partSize, routines = copyCommand.preparePartOption(MaxInt64)
    c.Assert(partSize, Equals, int64(922337203685478))
    c.Assert(routines, Equals, 10)

    p := 7 
    parallel := strconv.Itoa(p) 
    copyCommand.command.options[OptionParallel] = &parallel
    partSize, routines = copyCommand.preparePartOption(1)
    c.Assert(routines, Equals, p)
    str := ""
    copyCommand.command.options[OptionParallel] = &str
}

func (s *OssutilCommandSuite) TestResumeDownloadRetry(c *C) {
    bucketName := bucketNamePrefix + "cpnotexist"
    bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)

    err = copyCommand.ossResumeDownloadRetry(bucket, "", "", 0, 0)
    c.Assert(err, NotNil)
}

func (s *OssutilCommandSuite) TestCPIDKey(c *C) {
    bucket := bucketNameExist 

    object := "testobject" 

    ufile := "ossutil_test.cpidkey"
    data := "欢迎使用ossutil"
    s.createFile(ufile, data, c)

    cfile := "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", bucket, "abc", bucket, "abc") 
    s.createFile(cfile, data, c)

    command := "cp"
    str := ""
    args := []string{ufile, CloudURLToString(bucket, object)}
    ok := true
    routines := strconv.Itoa(Routines)
    thre := strconv.FormatInt(DefaultBigFileThreshold, 10)
    cpDir := CheckpointDir
    options := OptionMapType{
        "endpoint": &str,
        "accessKeyID": &str,
        "accessKeySecret": &str,
        "stsToken": &str,
        "configFile": &cfile,
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err := cm.RunCommand(command, args, options)
    c.Assert(err, NotNil)

    options = OptionMapType{
        "endpoint": &endpoint,
        "accessKeyID": &accessKeyID,
        "accessKeySecret": &accessKeySecret,
        "stsToken": &str,
        "configFile": &cfile,
        "force": &ok,
        "bigfileThreshold": &thre,
        "checkpointDir": &cpDir,
        "routines": &routines,
    }
    showElapse, err = cm.RunCommand(command, args, options)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)

    _ = os.Remove(ufile)
    _ = os.Remove(cfile)
}

func (s *OssutilCommandSuite) TestUploadOutputDir(c *C) {
    dir := "ossutil_test_output_dir" 
    _ = os.RemoveAll(dir)

    bucket := bucketNameExist
    object := randStr(10) 
    ufile := "ossutil_test.testoutputdir"
    data := "content" 
    s.createFile(ufile, data, c)

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // NoSuchBucket err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucketNameNotExist, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // SignatureDoesNotMatch err copy -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s", endpoint, accessKeyID, "abc", bucket, endpoint) 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucket, "abc") 
    s.createFile(configFile, data, c)

    // err copy without -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy with -r -> outputdir
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, dir) 
    os.Stdout = out
    str := s.readFile(resultPath, c)
    c.Assert(strings.Contains(str, "Error occurs, see more information in file"), Equals, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 

    _ = os.Remove(configFile)
    configFile = cfile

    // get file list of outputdir
    fileList, err := s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    // get report file content
    result := s.getReportResult(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileList[0]), c)
    c.Assert(len(result), Equals, 1)
    
    _ = os.Remove(ufile)
    _ = os.RemoveAll(dir)

    // err list with -C -> no outputdir
    udir := "TestUploadOutputDir"
    err = os.MkdirAll(udir, 0755)
    c.Assert(err, IsNil)
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucket, object), false, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.RemoveAll(udir)
}

func (s *OssutilCommandSuite) TestBatchUploadOutputDir(c *C) {
    udir := "TestBatchUploadOutputDir/" 
    _ = os.RemoveAll(udir)
    err := os.MkdirAll(udir, 0755)
    c.Assert(err, IsNil)

    num := 2 
    filePaths := []string{}
    for i := 0; i < num; i++ {
        filePath := randStr(10) 
        s.createFile(udir + "/" + filePath, fmt.Sprintf("测试文件：%d内容", i), c)
        filePaths = append(filePaths, filePath)
    }

    dir := "ossutil_test_output_dir" 
    _ = os.RemoveAll(dir)
    bucket := bucketNameExist

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(udir, CloudURLToString(bucket, udir + "/"), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy without -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucket, udir + "/"), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // NoSuchBucket err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucketNameNotExist, udir + "/"), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy -> outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n", "abc", accessKeyID, accessKeySecret) 
    s.createFile(configFile, data, c)

    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    showElapse, err = s.rawCPWithOutputDir(udir, CloudURLToString(bucket, udir + "/"), true, true, false, 1, dir) 
    os.Stdout = out
    str := s.readFile(resultPath, c)
    c.Assert(strings.Contains(str, "Error occurs, see more information in file"), Equals, true)
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 

    // get file list of outputdir
    fileList, err := s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    // get report file content
    result := s.getReportResult(fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileList[0]), c)
    c.Assert(len(result), Equals, num)
 
    _ = os.Remove(configFile)
    configFile = cfile

    _ = os.RemoveAll(udir)
    _ = os.RemoveAll(dir)
}

func (s *OssutilCommandSuite) TestDownloadOutputDir(c *C) {
    dir := "ossutil_test_output_dir" 
    _ = os.RemoveAll(dir)

    bucket := bucketNameExist
    object := randStr(10)
    s.putObject(bucket, object, uploadFileName, c)

    // normal copy without -r -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(bucket, object), downloadFileName, false, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // normal copy with -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucket, object), downloadDir, true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketNameNotExist, object), downloadFileName, true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy without -r -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucket, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucket, object), downloadFileName, false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // list err copy with -r -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucket, object), downloadDir, true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.RemoveAll(dir)
    _ = os.Remove(configFile)
    configFile = cfile
}

func (s *OssutilCommandSuite) TestCopyOutputDir(c *C) { 
    dir := "ossutil_test_output_dir" 
    _ = os.RemoveAll(dir)

    srcBucket := bucketNameExist
    destBucket := bucketNameDest

    object := randStr(10)
    s.putObject(srcBucket, object, uploadFileName, c)

    // normal copy -> no outputdir
    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(bucketNameNotExist, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(bucketNameNotExist, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // list err copy without -r -> no outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, srcBucket, "abc") 
    s.createFile(configFile, data, c)
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), false, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, object), CloudURLToString(destBucket, object), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir)
}

func (s *OssutilCommandSuite) TestBatchCopyOutputDir(c *C) {
    dir := "ossutil_test_output_dir" 
    _ = os.RemoveAll(dir)

    srcBucket := bucketNameExist
    destBucket := bucketNameDest

    objectList := []string{}
    num := 3
    for i := 0; i < num; i++ {
        object := randStr(10)
        s.putObject(srcBucket, object, uploadFileName, c)
        objectList = append(objectList, object)
    }

    showElapse, err := s.rawCPWithOutputDir(CloudURLToString(srcBucket, objectList[0]), CloudURLToString(destBucket, ""), true, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)

    _ = os.RemoveAll(dir)

    // normal copy -> no outputdir 
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), true, true, false, 1, dir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // bucketNameNotExist err copy -> no outputdir
    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(bucketNameNotExist, ""), true, true, false, 1, dir) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // test objectStatistic err
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data := fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", "abc", "def", "ghi", srcBucket, "abc", srcBucket, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(CloudURLToString(srcBucket, ""), CloudURLToString(destBucket, ""), true, true, false, 1, dir)
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir)
}

func (s *OssutilCommandSuite) TestConfigOutputDir(c *C) {
    // test default outputdir
    edir := "" 
    dir := "testoutputdir"
    dir1 := "newoutputdir"
    _ = os.RemoveAll(DefaultOutputDir)
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(dir1)

    bucket := bucketNameExist
    object := randStr(10) 
    ufile := "ossutil_test.testoutputdir"
    data := "content" 
    s.createFile(ufile, data, c)

    // err copy -> outputdir
    cfile := configFile
    configFile = "ossutil_test.config_boto"
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\n[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, bucket, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err := s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, edir) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, IsNil) 

    // get file list of outputdir
    fileList, err := s.getFileList(DefaultOutputDir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.RemoveAll(DefaultOutputDir) 

    // config outputdir
    data = fmt.Sprintf("[Credentials]\nendpoint=%s\naccessKeyID=%s\naccessKeySecret=%s\noutputDir=%s\n[Bucket-Endpoint]\n%s=%s[Bucket-Cname]\n%s=%s", endpoint, accessKeyID, accessKeySecret, dir, bucket, endpoint, bucket, "abc") 
    s.createFile(configFile, data, c)

    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, "") 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir)
    c.Assert(err, IsNil) 
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // get file list of outputdir
    fileList, err = s.getFileList(dir)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(DefaultOutputDir)

    // option and config
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, dir1) 
    c.Assert(err, IsNil)
    c.Assert(showElapse, Equals, true)
    _, err = os.Stat(dir1)
    c.Assert(err, IsNil) 
    _, err = os.Stat(dir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)
    _, err = os.Stat(DefaultOutputDir)
    c.Assert(err, NotNil) 
    c.Assert(os.IsNotExist(err), Equals, true)

    // get file list of outputdir
    fileList, err = s.getFileList(dir1)
    c.Assert(err, IsNil)
    c.Assert(len(fileList), Equals, 1)

    _ = os.Remove(configFile)
    configFile = cfile
    _ = os.RemoveAll(dir1)
    _ = os.RemoveAll(dir)
    _ = os.RemoveAll(DefaultOutputDir)

    s.createFile(uploadFileName, content, c)
    showElapse, err = s.rawCPWithOutputDir(ufile, CloudURLToString(bucket, object), true, true, false, 1, uploadFileName) 
    c.Assert(err, NotNil)
    c.Assert(showElapse, Equals, false)
}

func (s *OssutilCommandSuite) TestInitReportError(c *C) {
    s.createFile(uploadFileName, content, c)
    report, err := GetReporter(false, DefaultOutputDir, "")
    c.Assert(err, IsNil)
    c.Assert(report, IsNil)

    report, err = GetReporter(true, uploadFileName, "")
    c.Assert(err, NotNil)
    c.Assert(report, IsNil)
}

func (s *OssutilCommandSuite) TestCopyFunction(c *C) {
    // test fileStatistic
    copyCommand.monitor.init(operationTypePut)
    storageURL, err := StorageURLFromString("&~")
    c.Assert(err, IsNil)
    copyCommand.fileStatistic([]StorageURLer{storageURL})
    c.Assert(copyCommand.monitor.seekAheadEnd, Equals, true)
    c.Assert(copyCommand.monitor.seekAheadError, NotNil)

    // test fileProducer
    chFiles := make(chan fileInfoType, ChannelBuf)
    chListError := make(chan error, 1)
    storageURL, err = StorageURLFromString("&~")
    c.Assert(err, IsNil)
    copyCommand.fileProducer([]StorageURLer{storageURL}, chFiles, chListError)
    err = <- chListError
    c.Assert(err, NotNil)

    // test put object error
    bucketName := bucketNameNotExist
    bucket, err := copyCommand.command.ossBucket(bucketName)
    c.Assert(err, IsNil)
    err = copyCommand.ossPutObjectRetry(bucket, "object", "")
    c.Assert(err, NotNil)

    // test formatResultPrompt
    testResultFile, _ = os.OpenFile(resultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0664)
    out := os.Stdout
    os.Stdout = testResultFile
    err = fmt.Errorf("test error")
    copyCommand.cpOption.ctnu = true
    err = copyCommand.formatResultPrompt(err)
    c.Assert(err, IsNil)
    os.Stdout = out
    str := s.readFile(resultPath, c)
    c.Assert(strings.Contains(str, "Error"), Equals, true)

    // test download file error
    err = copyCommand.ossDownloadFileRetry(bucket, "object", downloadFileName)
    c.Assert(err, NotNil)

    // test truncateFile
    err = copyCommand.truncateFile("ossutil_notexistfile", 1)
    c.Assert(err, NotNil)
    s.createFile(uploadFileName, "abc", c)
    err = copyCommand.truncateFile(uploadFileName, 1)
    c.Assert(err, IsNil)
    str = s.readFile(uploadFileName, c)
    c.Assert(str, Equals, "a")
}
