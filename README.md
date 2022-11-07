# flow

flow is high throughput, low-latency service that provides an event-type agnostic API to ingest clickstream data from client apps such as mobile apps, and sites. The service collects, processes, enriches, and stores the clickstreams for real-time analytics to get a 360-view of the user for personalization.

This architecture shows how you can collect, process, enriches, and store clickstream at scale using Azure Functions, Azure CosmosDB, Azure Event Hub, and Azure Blob Store.

## Architecture

- [Event Hubs](https://azure.microsoft.com/services/event-hubs) ingests raw click-stream data from Azure Functions and  archives raw click-stream data to Azure Storage.
- [Azure Cosmos DB](https://azure.microsoft.com/services/cosmos-db) stores aggregated data of clicks by user, product, and offers user-profile information.
- [Azure Storage](https://azure.microsoft.com/services/storage) stores archived raw click-stream data from Event Hubs.
- [Power BI](https://powerbi.microsoft.com/) enables visualization of user activity data and offers presented by reading in data from Azure Cosmos DB.

