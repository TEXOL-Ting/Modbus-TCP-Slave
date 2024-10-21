package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

// var ModbusPort = 502 // Default port for Modbus TCP

var Config_global Config

// Simulated holding registers and input registers (data storage for the Modbus slave)
// var holdingRegisters = make([]uint16, 10)
var inputRegisters = make([]uint16, 200) // Input registers

func handleModbusRequest(conn net.Conn) {
	defer conn.Close()

	for {
		// Receive the request from the master
		request := make([]byte, 12) // A Modbus TCP request typically has 12 bytes
		n, err := conn.Read(request)
		if err != nil {
			log.Printf("Error reading request: %v", err)
			return
		}

		println("request:")
		display(request)

		// Ensure we have received at least 12 bytes (Modbus standard request length)
		if n < 12 {
			log.Printf("Received incomplete request: %d bytes", n)
			return
		}

		// Extract the function code (7th byte in the request)
		functionCode := request[7]

		// Check which function code is used
		switch functionCode {
		/*
			case 0x03: // Read Holding Registers
			handleReadHoldingRegisters(conn, request)
		*/
		case 0x04: // Read Input Registers
			handleReadInputRegisters(conn, request)
		default:
			log.Printf("Unsupported function code: %x", functionCode)
		}
	}
}

/*
// Handle the function code 0x03 (Read Holding Registers)
func handleReadHoldingRegisters(conn net.Conn, request []byte) {
	// Extract the starting address and quantity of registers from the request
	startAddress := binary.BigEndian.Uint16(request[8:10])
	quantity := binary.BigEndian.Uint16(request[10:12])

	// Check if the request is valid
	if startAddress+quantity > uint16(len(holdingRegisters)) {
		log.Println("Requested range exceeds available holding registers")
		return
	}

	// Build the response
	response := make([]byte, 9+quantity*2)
	copy(response[:6], request[:6])    // Copy the transaction ID, protocol ID, etc.
	response[6] = 0                    // Length high byte
	response[7] = byte(3 + quantity*2) // Length low byte
	response[8] = request[7]           // Function code
	response[9] = byte(quantity * 2)   // Number of bytes

	// Write the values of the holding registers into the response
	for i := 0; i < int(quantity); i++ {
		binary.BigEndian.PutUint16(response[9+i*2:], holdingRegisters[startAddress+uint16(i)])
	}

	// Send the response to the client
	_, err := conn.Write(response)
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	log.Printf("Handled Read Holding Registers request, address: %d, quantity: %d", startAddress, quantity)
}
*/

// Handle the function code 0x04 (Read Input Registers)
func handleReadInputRegisters(conn net.Conn, request []byte) {
	// Extract the starting address and quantity of registers from the request
	startAddress := binary.BigEndian.Uint16(request[8:10])
	quantity := binary.BigEndian.Uint16(request[10:12])

	println("startAddress", startAddress, "quantity", quantity)

	// Check if the request is valid
	if startAddress+quantity > uint16(len(inputRegisters)) {
		log.Println("Requested range exceeds available input registers")
		return
	}

	// Build the response
	response := make([]byte, 9+quantity*2)
	copy(response[:4], request[:4])    // Copy the transaction ID, protocol ID
	response[4] = 0                    // Length high byte
	response[5] = byte(3 + quantity*2) // Length low byte
	response[6] = request[6]           // Unit ID
	response[7] = request[7]           // Function code
	response[8] = byte(quantity * 2)   // Number of bytes

	//println("response1:")
	//display(response)

	// Write the values of the input registers into the response
	for i := 0; i < int(quantity); i++ {
		binary.BigEndian.PutUint16(response[9+i*2:], inputRegisters[startAddress+uint16(i)])
	}

	println("response:")
	display(response)

	// Send the response to the client
	_, err := conn.Write(response)
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	log.Printf("Handled Read Input Registers request, address: %d, quantity: %d", startAddress, quantity)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./Simulation_Modbus <filename>")
		return
	}

	filename := os.Args[1]

	// Read the config back from the YAML file
	readConfig, err := readConfig(filename)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	fmt.Printf("Config read from file '%s': %+v\n", filename, readConfig)

	Config_global = readConfig

	/*
		// Initialize holding registers and input registers with some data
		for i := range holdingRegisters {
			holdingRegisters[i] = uint16(i + 1)
		}

			for i := range inputRegisters {
				inputRegisters[i] = uint16(i + 1)
			}
	*/

	go Simulation()

	// Start listening on the TCP port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", Config_global.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	defer listener.Close()
	log.Printf("Modbus TCP slave started, listening on port %d", Config_global.Port)

	for {
		// Accept client connections
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		// Handle the client request
		go handleModbusRequest(conn)
	}
}

func Simulation() {

	var feature_range [][]float32

	//oa_range := [][]float32{{2.6, 2.9}, {2.8, 3.1}, {4.4, 4.8}}
	feature_range = append(feature_range, get_range(Config_global.OA_Axial))
	feature_range = append(feature_range, get_range(Config_global.OA_Radial_1))
	feature_range = append(feature_range, get_range(Config_global.OA_Radial_2))

	peak_range := [][]float32{{6.6, 7.9}, {6.6, 7.9}, {10.3, 14.4}}
	feature_range = append(feature_range, peak_range...)

	var ptop_range = make([][]float32, len(peak_range))
	for i, ra := range peak_range {
		var x = make([]float32, len(ra))
		for j, rm := range ra {
			x[j] = rm * 2
		}
		ptop_range[i] = x
	}
	feature_range = append(feature_range, ptop_range...)

	cf_range := [][]float32{{2.6, 3.1}, {2.6, 3.1}, {2.6, 3.1}}
	feature_range = append(feature_range, cf_range...)

	var pib_range = make([][]float32, 10)
	for i := range pib_range {
		pib_range[i] = []float32{float32(i) * 0.1, float32(i+1) * 0.1}
	}
	feature_range = append(feature_range, pib_range...) //X
	feature_range = append(feature_range, pib_range...) //Y
	feature_range = append(feature_range, pib_range...) //Z

	failure_range := [][]float32{{7.3, 7.5}, {8.2, 8.6}, {2.1, 2.3}, {0.3, 0.4}, {0, 0}, {0, 0}}
	failure_range[0] = get_range(Config_global.Def_Bea)
	failure_range[3] = get_range(Config_global.Def_Imb)
	feature_range = append(feature_range, failure_range...)
	feature_range = append(feature_range, failure_range...) //%

	var SensorID = asciiToUint16(Config_global.SensorID) //[]uint16{0x3033, 0x3030, 0x3030}
	copy(inputRegisters[4:], SensorID)

	inputRegisters[7] = 1860                            //RPM
	copy(inputRegisters[116:], u32ToUint16Slice(27200)) //Temperature
	inputRegisters[167] = 1                             //Slave ID
	inputRegisters[168] = 0                             //Machine type

	rand.Seed(time.Now().UnixNano())

	var count uint16 = 0
	ticker := time.NewTicker(1 * time.Second) // 创建一个每秒触发的 ticker
	defer ticker.Stop()                       // 确保程序退出前停止 ticker
	for range ticker.C {                      // 每次 ticker 触发时都会进入此循环
		data(&count, feature_range) // 执行函数
	}
}

func data(count *uint16, feature_range [][]float32) {
	fmt.Println("Task executed at", time.Now())

	var feature_list []uint32

	*count += 1
	var Output_Count = make([]uint16, 4)
	Output_Count[3] = *count
	copy(inputRegisters, Output_Count) //Data Output Count

	for i, r := range feature_range {
		Address := i*2 + 8
		feature := random_feature(r)
		copy(inputRegisters[Address:], u32ToUint16Slice(feature))
		feature_list = append(feature_list, feature)
	}

	for i := 0; i < 3; i++ {
		var qc uint16 = 0
		oa := feature_list[i]
		if oa < uint32(Config_global.OA_QC[0]*1000) {
			qc = 0
		} else if oa < uint32(Config_global.OA_QC[1]*1000) {
			qc = 1
		} else {
			qc = 2
		}
		inputRegisters[118+i] = qc
		//fmt.Println(i, "OA: ", oa, "QC: ", qc)
	}

}

func random_feature(range_list []float32) uint32 {
	//var oa_range = [][]float32{{2.6, 2.9}, {2.8, 3.1}, {4.4, 4.8}}
	var feature uint32

	min := int(range_list[0] * 1000)
	max := int(range_list[1] * 1000)
	feature = uint32(rand.Intn(max-min+1) + min)
	//println("min:", min, "max:", max, "value:", feature)

	return feature
}

func u32ToUint16Slice(value uint32) []uint16 {
	// 将 uint32 拆分为两个 uint16，低 16 位和高 16 位
	low := uint16(value & 0xFFFF)          // 低 16 位
	high := uint16((value >> 16) & 0xFFFF) // 高 16 位

	return []uint16{high, low}
}

func asciiToUint16(s string) []uint16 {
	result := make([]uint16, 0, len(s)/2)

	for i := 0; i < 3; i++ {
		result = append(result, uint16(s[i*2])<<8|uint16(s[i*2+1]))
	}

	return result
}

func display(r []byte) {
	var cmd string
	for _, b := range r {
		cmd = cmd + fmt.Sprintf("%02X", b) + " "
	}

	println(cmd)
}

func get_range(r []float32) []float32 {
	return []float32{r[0] - r[1], r[0] + r[1]}
}
