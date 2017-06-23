
//包名
package main

//导入包
import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("第一个chiancode-01")

// 链码的实现
type SimpleChaincode struct {
}

//初始化函数
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	logger.Info("########### 函数初始化-02 ###########")

	_, args := stub.GetFunctionAndParameters() //下划线意思是忽略这个变量
	var A, B string    
	var Aval, Bval int 
	var err error

	if len(args) != 4 {
		return shim.Error("输入的参数不正确！")
	}

	// 初始化链码
	A = args[0]
	Aval, err = strconv.Atoi(args[1]) //将字符串转换为十进制整数
	if err != nil {
		return shim.Error("请输入数值")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3]) //将字符串转换为十进制整数
	if err != nil {
		return shim.Error("请输入数值")
	}
	logger.Info("打印ab的值："+"Aval = %d, Bval = %d\n", Aval, Bval)

	//把状态写在分类帐上
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	//把状态写在分类帐上
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)


}

// 交易
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### 执行 invoke 方法 ###########")

	function, args := stub.GetFunctionAndParameters() //获取参数名和参数值
	
	if function != "invoke" {
		return shim.Error("函数输入不正确，请输入invoke")
	}

	if len(args) < 2 {
		return shim.Error("参数至少输入2个")
	}

	if args[0] == "delete" {
		// 从状态删除一个实体
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// 从状态查询一个实体
		return t.query(stub, args)
	}
	if args[0] == "move" {
		// a向b 转账
		return t.move(stub, args)
	}

	logger.Errorf("第一个参数可能是delete或query或move", args[0])
	return shim.Error(fmt.Sprintf("请检查您的参数输入是否正确！", args[0]))
}

//转账
func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 转账
	var A, B string    
	var Aval, Bval int 
	var X int          // 每次交易值
	var err error

	if len(args) != 4 {
		return shim.Error("参数输入不正确")
	}

	A = args[1]
	B = args[2]

	// 从分类帐中获取状态
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("没有查到")
	}
	if Avalbytes == nil {
		return shim.Error("没有找到该参数")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes)) //将字符串转换为十进制整数

	// 从分类帐中获取状态
	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("没有查到")
	}
	if Bvalbytes == nil {
		return shim.Error("没有找到该参数")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))  //将字符串转换为十进制整数

	// 获取交易的值
	X, err = strconv.Atoi(args[3]) //将字符串转换为十进制整数
	if err != nil {
		return shim.Error("请输入数值")
	}
	Aval = Aval - X
	Bval = Bval + X
	logger.Infof("a b 交易后的值："+"Aval = %d, Bval = %d\n", Aval, Bval)

	// 写进分类账本
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

    // 写进分类账本
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

// 删除
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数输入不正确")
	}

	A := args[1]

	// 从分类账本中删除
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("删除失败！")
	}

	return shim.Success(nil)
}

// 查询
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var A string 
	var err error

	if len(args) != 2 {
		return shim.Error("参数输入不正确")
	}

	A = args[1]

	// 从分类账本中获取
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"错误\":\"获取失败 " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"错误\":\"没有这个值 " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"交易名\":\"" + A + "\",\"交易值\":\"" + string(Avalbytes) + "\"}"
	logger.Infof("查询结果:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("chaincode启动失败: %s", err)
	}
}
