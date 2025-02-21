package examples

func (o *OptionForTest) SetNameMessageType(v string) error {
	o.nameMessageType = v

	return nil
}

func (o *OptionForTest) SetMaxLogFileSize(v int) error {
	o.maxLogFileSize = v

	return nil
}

func (o *OptionForTest) SetPathDirectory(v string) error {
	o.pathDirectory = v

	return nil
}

func (o *OptionForTest) SetWritingStdout(v bool) {
	o.writingStdout = v
}

func (o *OptionForTest) SetWritingFile(v bool) {
	o.writingFile = v
}

func (o *OptionForTest) SetWritingDB(v bool) {
	o.writingDB = v
}

func (o *OptionForTest) GetNameMessageType() string {
	return o.nameMessageType
}

func (o *OptionForTest) GetMaxLogFileSize() int {
	return o.maxLogFileSize
}

func (o *OptionForTest) GetPathDirectory() string {
	return o.pathDirectory
}

func (o *OptionForTest) GetWritingStdout() bool {
	return o.writingStdout
}

func (o *OptionForTest) GetWritingFile() bool {
	return o.writingFile
}

func (o *OptionForTest) GetWritingDB() bool {
	return o.writingDB
}
