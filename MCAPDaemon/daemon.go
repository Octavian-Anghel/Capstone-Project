package main

import (
	"bytes"
	//"context"
	"crypto/x509"
	"encoding/json"
	"path"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/fsnotify/fsnotify"
	"github.com/hyperledger/fabric-gateway/pkg/client"                                 
	"github.com/hyperledger/fabric-gateway/pkg/hash"                                   
	"github.com/hyperledger/fabric-gateway/pkg/identity"                               
	//"github.com/hyperledger/fabric-protos-go-apiv2/gateway"                            
	"google.golang.org/grpc"                                                           
	"google.golang.org/grpc/credentials"                                               
	//"google.golang.org/grpc/status"

)

const (                                                                                    
    mspID        = "Org1MSP"                                                               
    cryptoPath   = "/home/oz/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"   
    certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"              
    keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"               
    tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"                 
    peerEndpoint = "dns:///localhost:7051"                                                 
    gatewayPeer  = "peer0.org1.example.com"                                                
)   

var now = time.Now()                                                                       
var assetId = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)

// newGrpcConnection creates a gRPC connection to the Gateway server.                      
func newGrpcConnection() *grpc.ClientConn {                                                
    certificatePEM, err := os.ReadFile(tlsCertPath)                                        
    if err != nil {                                                                        
        panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))                   
    }                                                                                      
                                                                       
    certificate, err := identity.CertificateFromPEM(certificatePEM)                        
    if err != nil {                                                                        
        panic(err)                                                                         
    }                                                                                      
                                                                                           
    certPool := x509.NewCertPool()                                                         
    certPool.AddCert(certificate)                                                          
    transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)        
                                                                                           
    connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
    if err != nil {                                                                        
        panic(fmt.Errorf("failed to create gRPC connection: %w", err))                     
    }                                                                                      
                                                                                           
    return connection                                                                      
}                                                                                          
                                                                                           
// newIdentity creates a client identity for this Gateway connection using an X.509 certifi
func newIdentity() *identity.X509Identity {                                                
    certificatePEM, err := readFirstFile(certPath)                                         
    if err != nil {                                                                        
        panic(fmt.Errorf("failed to read certificate file: %w", err))                      
    }                                                                                      
                                                                                           
    certificate, err := identity.CertificateFromPEM(certificatePEM)                        
    if err != nil {                                                                        
        panic(err)                                                                         
    }                                                                                      
                                                                                           
    id, err := identity.NewX509Identity(mspID, certificate)                                
    if err != nil {                                                                        
        panic(err)                                                                         
    }                                                                                      
                                                                                           
    return id                                                                              
}

// newSign creates a function that generates a digital signature from a message digest usin
func newSign() identity.Sign {                                                             
    privateKeyPEM, err := readFirstFile(keyPath)                                           
    if err != nil {                                                                        
        panic(fmt.Errorf("failed to read private key file: %w", err))                      
    }                                                                                      
                                                                                           
    privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)                           
    if err != nil {                                                                        
        panic(err)                                                                         
    }                                                                                      
                                                                                           
    sign, err := identity.NewPrivateKeySign(privateKey)                                    
    if err != nil {                                                                        
        panic(err)                                                                         
    }                                                                                      
                                                                                           
    return sign                                                                            
}                                                                                          
                                                                                           
func readFirstFile(dirPath string) ([]byte, error) {                                       
    dir, err := os.Open(dirPath)                                                           
    if err != nil {                                                                        
            return nil, err                                                                    
        }                                                                                      
                                                                                               
        fileNames, err := dir.Readdirnames(1)                                                  
        if err != nil {                                                                        
            return nil, err                                                                    
        }                                                                                      
                                                                                               
        return os.ReadFile(path.Join(dirPath, fileNames[0]))                                   
    }

// Submit a transaction synchronously, blocking until it has been committed to the ledger.                                                                        
func CreateAsset(contract *client.Contract, hash string, mcapID string, operationID string) {                                                                                                                     
    fmt.Printf("\n--> Submit Transaction: CreateAsset, creates new hash asset on the ledger\n")                            

    _, err := contract.SubmitTransaction("CreateAsset", string(time.Now().Format(time.RFC3339)), hash, mcapID, operationID)                                                                    
    if err != nil {                                                                                                                                               
        fmt.Printf("failed to submit transaction: %w", err)                                                                                               
    }                                                                                                                                                             
                                                                                                                                                                          
    fmt.Printf("*** Transaction committed successfully\n")                                                                                                        
}

func printTime(format string, a ...interface{}) {
		fmt.Printf("[%s] ", time.Now().Format(time.RFC3339))
	fmt.Printf(format+"\n", a...)
}

func exit(format string, a ...interface{}) {
		printTime(format, a...)
	os.Exit(1)
}

func getMagicBytes(filePath string) (string, error) {
		file, err := os.Open(filePath)
	if err != nil {
			return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
			return "", err
	}
	fileSize := fileInfo.Size()

	if fileSize < 5 {
			return "", errors.New("filesize too small, likely bad file")
	}

	start := fileSize - int64(7)
	buff := make([]byte, 7)

	_, err = file.ReadAt(buff, start)
	if err != nil {
			return "", err
	}

	return string(buff), nil
}

func dedupLoop(w *fsnotify.Watcher, contract *client.Contract) {
		var (
		waitFor    = 100 * time.Millisecond
		mu         sync.Mutex
		timers     = make(map[string]*time.Timer)
		printEvent = func(e fsnotify.Event) {
				printTime("Detected event: %s", e)

			if strings.HasSuffix(e.Name, ".mcap") {
					magic, err := getMagicBytes(e.Name)
				if err != nil {
						fmt.Printf("Failed to read magic bytes from %s: %v\n", e.Name, err)
				} else {
						fmt.Printf("Magic bytes from %s: %s\n", e.Name, magic)
					if magic == "MCAP0\r\n" {
							fmt.Println("Valid MCAP file detected! pushing over to the hash and upload daemon\n")
						// Additional processing here
							hashString := HashNUpload(e.Name)
							CreateAsset(contract, hashString, e.Name, "Dummy Operation ID")
					}
				}
			}

			mu.Lock()
			delete(timers, e.Name)
			mu.Unlock()
		}
	)

	for {
			select {
			case err, ok := <-w.Errors:
			if !ok {
					return
			}
			printTime("ERROR: %s", err)

		case e, ok := <-w.Events:
			if !ok {
					return
			}

			if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
					continue
			}

			mu.Lock()
			t, ok := timers[e.Name]
			mu.Unlock()

			if !ok {
					t = time.AfterFunc(math.MaxInt64, func() { printEvent(e) })
				t.Stop()
				mu.Lock()
				timers[e.Name] = t
				mu.Unlock()
			}

			t.Reset(waitFor)
		}
	}
}

// Format JSON data                                                                        
func formatJSON(data []byte) string {                                                      
    var prettyJSON bytes.Buffer                                                            
    if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {                       
        panic(fmt.Errorf("failed to parse JSON: %w", err))                                 
    }                                                                                      
    return prettyJSON.String()                                                             
}            

func main() {

	clientConnection := newGrpcConnection()                                                
	defer clientConnection.Close()                                                         
	                                                                                           
	id := newIdentity()                                                                    
	sign := newSign()                                                                      
	                                                                                           
    // Create a Gateway connection for a specific client identity                          
    gw, err := client.Connect(                                                             
	   id,                                                                                
	   client.WithSign(sign),                                                             
	   client.WithHash(hash.SHA256),                                                      
	   client.WithClientConnection(clientConnection),                                     
	   // Default timeouts for different gRPC calls                                       
	   client.WithEvaluateTimeout(5*time.Second),                                         
	   client.WithEndorseTimeout(15*time.Second),                                         
	   client.WithSubmitTimeout(5*time.Second),                                           
	   client.WithCommitStatusTimeout(1*time.Minute),                                     
    )                                                                                      
    if err != nil {                                                                      
        panic(err)                                                                         
    }                                                                                
    defer gw.Close()

	chaincodeName := "mcap"
	channelName := "mychannel"

	network := gw.GetNetwork(channelName)                                                  
	contract := network.GetContract(chaincodeName) 

	w, err := fsnotify.NewWatcher()
	if err != nil {
		exit("creating a new watcher: %s", err)
	}
	defer w.Close()

	go dedupLoop(w, contract)

	path := "/shared"
	err = w.Add(path)
	if err != nil {
		exit("%q: %s", path, err)
	}

	// Prevent main from exiting
	select {}
}
