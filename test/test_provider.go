package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	devopsProvider "github.com/hashicorp/terraform-provider-scaffolding-framework/internal/provider"
)

func main() {
	// Create a new provider instance
	p := devopsProvider.New("dev")()

	// Test provider metadata
	ctx := context.Background()
	
	// Test metadata
	var metadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &metadataResp)
	
	fmt.Printf("Provider Type Name: %s\n", metadataResp.TypeName)
	fmt.Printf("Provider Version: %s\n", metadataResp.Version)
	
	// Test schema
	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)
	
	if schemaResp.Diagnostics.HasError() {
		log.Fatal("Schema has errors:", schemaResp.Diagnostics.Errors())
	}
	
	fmt.Println("✅ Provider schema is valid!")
	
	// Test data sources
	dataSources := p.DataSources(ctx)
	fmt.Printf("📊 Data Sources available: %d\n", len(dataSources))
	
	// Test resources  
	resources := p.Resources(ctx)
	fmt.Printf("🔧 Resources available: %d\n", len(resources))
	
	fmt.Println("\n🎉 Provider validation successful!")
}
