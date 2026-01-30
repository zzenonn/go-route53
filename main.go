package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func main() {
	// Load AWS configuration using default credential chain
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create Route53 client
	client := route53.NewFromConfig(cfg)

	// Example: Get all zones and their records
	getAllZonesAndRecords(client)

	// Example: Create a subdomain
	// Uncomment to test:
	// err = createSubdomain(client, "Z1234567890ABC", "test123.example.com", "192.0.2.1", 300)
	// if err != nil {
	// 	log.Printf("failed to create subdomain: %v", err)
	// } else {
	// 	fmt.Println("Subdomain created successfully")
	// }

	// Example: Delete a subdomain
	// Uncomment to test:
	// err = deleteSubdomain(client, "Z1234567890ABC", "test123.example.com", "192.0.2.1", 300)
	// if err != nil {
	// 	log.Printf("failed to delete subdomain: %v", err)
	// } else {
	// 	fmt.Println("Subdomain deleted successfully")
	// }
}

// getAllZonesAndRecords retrieves and displays all hosted zones and their DNS records
func getAllZonesAndRecords(client *route53.Client) {
	zones, err := listHostedZones(client)
	if err != nil {
		log.Fatalf("failed to list hosted zones, %v", err)
	}

	fmt.Printf("Found %d hosted zones\n\n", len(zones))

	for _, zone := range zones {
		displayZoneRecords(client, zone)
	}
}

// displayZoneRecords displays all DNS records for a specific hosted zone
func displayZoneRecords(client *route53.Client, zone types.HostedZone) {
	fmt.Printf("Hosted Zone: %s (ID: %s)\n", *zone.Name, *zone.Id)
	fmt.Println("----------------------------------------")

	records, err := listResourceRecordSets(client, *zone.Id)
	if err != nil {
		log.Printf("failed to list records for zone %s: %v", *zone.Id, err)
		return
	}

	for _, record := range records {
		displayRecord(record)
	}
	fmt.Println()
}

// displayRecord prints details of a single DNS record
func displayRecord(record types.ResourceRecordSet) {
	fmt.Printf("  Name: %s\n", *record.Name)
	fmt.Printf("  Type: %s\n", record.Type)
	fmt.Printf("  TTL: %d\n", getRecordTTL(record))

	for _, rr := range record.ResourceRecords {
		fmt.Printf("    Value: %s\n", *rr.Value)
	}

	if record.AliasTarget != nil {
		fmt.Printf("    Alias Target: %s\n", *record.AliasTarget.DNSName)
	}

	fmt.Println()
}

func listHostedZones(client *route53.Client) ([]types.HostedZone, error) {
	var zones []types.HostedZone
	var marker *string

	for {
		input := &route53.ListHostedZonesInput{
			Marker: marker,
		}

		result, err := client.ListHostedZones(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		zones = append(zones, result.HostedZones...)

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}

	return zones, nil
}

func listResourceRecordSets(client *route53.Client, zoneID string) ([]types.ResourceRecordSet, error) {
	var records []types.ResourceRecordSet
	var startRecordName *string
	var startRecordType types.RRType

	for {
		input := &route53.ListResourceRecordSetsInput{
			HostedZoneId:    &zoneID,
			StartRecordName: startRecordName,
			StartRecordType: startRecordType,
		}

		result, err := client.ListResourceRecordSets(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		records = append(records, result.ResourceRecordSets...)

		if !result.IsTruncated {
			break
		}
		startRecordName = result.NextRecordName
		startRecordType = result.NextRecordType
	}

	return records, nil
}

func getRecordTTL(record types.ResourceRecordSet) int64 {
	if record.TTL != nil {
		return *record.TTL
	}
	return 0
}

// createSubdomain creates a new A record subdomain in the specified hosted zone
func createSubdomain(client *route53.Client, hostedZoneID, subdomain, ipAddress string, ttl int64) error {
	changeBatch := &types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: types.ChangeActionCreate,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name: &subdomain,
					Type: types.RRTypeA,
					TTL:  &ttl,
					ResourceRecords: []types.ResourceRecord{
						{
							Value: &ipAddress,
						},
					},
				},
			},
		},
	}

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &hostedZoneID,
		ChangeBatch:  changeBatch,
	}

	_, err := client.ChangeResourceRecordSets(context.TODO(), input)
	return err
}

// deleteSubdomain deletes an A record subdomain from the specified hosted zone
func deleteSubdomain(client *route53.Client, hostedZoneID, subdomain, ipAddress string, ttl int64) error {
	changeBatch := &types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: types.ChangeActionDelete,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name: &subdomain,
					Type: types.RRTypeA,
					TTL:  &ttl,
					ResourceRecords: []types.ResourceRecord{
						{
							Value: &ipAddress,
						},
					},
				},
			},
		},
	}

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &hostedZoneID,
		ChangeBatch:  changeBatch,
	}

	_, err := client.ChangeResourceRecordSets(context.TODO(), input)
	return err
}
