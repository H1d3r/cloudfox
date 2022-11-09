package azure

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/BishopFox/cloudfox/constants"
	"github.com/BishopFox/cloudfox/utils"
	"github.com/aws/smithy-go/ptr"
)

func TestGetComputeRelevantData(t *testing.T) {
	fmt.Println()
	fmt.Println("[test case] GetComputeRelevantData Function")

	// Mocked functions to simulate Azure responses
	GetComputeVMsPerResourceGroup = func(subscriptionID, resourceGroup string) []compute.VirtualMachine {
		testFile, err := os.ReadFile(constants.VMS_TEST_FILE)
		if err != nil {
			log.Fatalf("could not read file %s", constants.VMS_TEST_FILE)
		}
		var vms []compute.VirtualMachine
		err = json.Unmarshal(testFile, &vms)
		if err != nil {
			log.Fatalf("could not unmarshall file %s", constants.VMS_TEST_FILE)
		}

		if subscriptionID == "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAA" || resourceGroup == "RG1" {
			return vms
		}
		return []compute.VirtualMachine{}
	}
	GetNICdetails = func(subscriptionID, resourceGroup string, nicReference compute.NetworkInterfaceReference) (network.Interface, error) {
		testFile, err := os.ReadFile(constants.NICS_TEST_FILE)
		if err != nil {
			log.Fatalf("could not read file %s", constants.NICS_TEST_FILE)
		}
		var nics []network.Interface
		err = json.Unmarshal(testFile, &nics)
		if err != nil {
			log.Fatalf("could not unmarshall file %s", constants.VMS_TEST_FILE)
		}
		nicName := strings.Split(ptr.ToString(nicReference.ID), "/")[len(strings.Split(ptr.ToString(nicReference.ID), "/"))-1]
		switch nicName {
		case "NetworkInterface1":
			return nics[0], nil
		case "NetworkInterface2":
			return nics[1], nil
		case "NetworkInterface3":
			return nics[2], nil
		case "NetworkInterface4":
			return nics[3], nil
		case "NetworkInterface5":
			return nics[4], nil
		default:
			return network.Interface{}, fmt.Errorf("nic not found: %s", ptr.ToString(nicReference.ID))
		}
	}
	GetPublicIPM = func(subscriptionID, resourceGroup string, ip network.InterfaceIPConfiguration) (*string, error) {
		publicIPID := ptr.ToString(ip.InterfaceIPConfigurationPropertiesFormat.PublicIPAddress.ID)
		publicIPName := strings.Split(publicIPID, "/")[len(strings.Split(publicIPID, "/"))-1]
		switch publicIPName {
		case "PublicIpAddress1A":
			return ptr.String("72.88.100.1"), nil
		case "PublicIpAddress1B":
			return ptr.String("72.88.100.2"), nil
		case "PublicIpAddress2A":
			return ptr.String("72.88.100.3"), nil
		case "PublicIpAddress3A":
			return ptr.String("72.88.100.3"), nil
		case "PublicIpAddress4A":
			return ptr.String("72.88.100.4"), nil
		case "PublicIpAddress5A":
			return ptr.String("72.88.100.5"), nil
		default:
			return nil, fmt.Errorf("public IP not found %s", publicIPName)
		}
	}

	// Test case parameters
	subscription := "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAA"
	resourceGroup := "RG1"
	header, body := GetComputeRelevantData(subscription, resourceGroup)
	verbosity := 2
	outputType := "table"
	outputPath := filepath.Join(constants.CLOUDFOX_BASE_OUTPUT_DIRECTORY, fmt.Sprintf("%s_%s", constants.AZ_OUTPUT_DIRECTORY, resourceGroup))
	fileName := constants.AZ_INTANCES_MODULE_NAME
	outputPrefixIdentifier := resourceGroup

	utils.MockFileSystem(true)
	utils.OutputSelector(verbosity, outputType, header, body, outputPath, fileName, constants.AZ_INTANCES_MODULE_NAME, outputPrefixIdentifier)
}
