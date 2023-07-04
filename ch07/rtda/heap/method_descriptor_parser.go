package heap

import "strings"

type MethodDescriptorParser struct {
	raw    string
	offset int
	parsed *MethodDescriptor
}

func parseMethodDescriptor(descriptor string) *MethodDescriptor {
	parser := &MethodDescriptorParser{}
	return parser.parse(descriptor)
}

func (self *MethodDescriptorParser) parse(descriptor string) *MethodDescriptor {
	self.raw = descriptor
	self.parsed = &MethodDescriptor{}
	self.startParams()
	self.parseParamTypes()
	self.endParams()
	self.parseReturnType()
	self.finish()
	return self.parsed
}

// 方法描述符前面为参数类型，最后一个是返回值类型，如 (II)V，所以方法描述符第一个符号一定是 '('
func (self *MethodDescriptorParser) startParams() {
	if self.readUint8() != '(' {
		self.causePanic()
	}
}

// 方法描述符前面为参数类型，最后一个是返回值类型，如 (II)V，所以方法描述符参数最后一个符号一定是 ')'
func (self *MethodDescriptorParser) endParams() {
	if self.readUint8() != ')' {
		self.causePanic()
	}
}
func (self *MethodDescriptorParser) finish() {
	if self.offset != len(self.raw) {
		self.causePanic()
	}
}

func (self *MethodDescriptorParser) causePanic() {
	panic("BAD descriptor: " + self.raw)
}

func (self *MethodDescriptorParser) readUint8() uint8 {
	b := self.raw[self.offset]
	self.offset++
	return b
}
func (self *MethodDescriptorParser) unreadUint8() {
	self.offset--
}

func (self *MethodDescriptorParser) parseParamTypes() {
	for {
		t := self.parseFieldType()
		if t != "" {
			self.parsed.addParameterType(t)
		} else {
			break
		}
	}
}

// 方法描述符前面为参数类型，最后一个是返回值类型，如 (II)V，所以方法描述符最后一个字符一定是方法返回值类型
func (self *MethodDescriptorParser) parseReturnType() {
	// 返回值类型是 void 类型
	if self.readUint8() == 'V' {
		self.parsed.returnType = "V"
		return
	}

	self.unreadUint8()
	t := self.parseFieldType()
	if t != "" {
		self.parsed.returnType = t
		return
	}

	self.causePanic()
}

// 解析参数类型
func (self *MethodDescriptorParser) parseFieldType() string {
	switch self.readUint8() {
	case 'B': // byte
		return "B"
	case 'C': // char
		return "C"
	case 'D': // double
		return "D"
	case 'F': // float
		return "F"
	case 'I': // int
		return "I"
	case 'J': // long
		return "J"
	case 'S': // string
		return "S"
	case 'Z': // boolean
		return "Z"
	case 'L': // long
		return self.parseObjectType()
	case '[': // 数组
		return self.parseArrayType()
	default:
		self.unreadUint8()
		return ""
	}
}

func (self *MethodDescriptorParser) parseObjectType() string {
	unread := self.raw[self.offset:]
	semicolonIndex := strings.IndexRune(unread, ';')
	if semicolonIndex == -1 {
		self.causePanic()
		return ""
	} else {
		objStart := self.offset - 1
		objEnd := self.offset + semicolonIndex + 1
		self.offset = objEnd
		descriptor := self.raw[objStart:objEnd]
		return descriptor
	}
}

func (self *MethodDescriptorParser) parseArrayType() string {
	arrStart := self.offset - 1
	self.parseFieldType()
	arrEnd := self.offset
	descriptor := self.raw[arrStart:arrEnd]
	return descriptor
}
