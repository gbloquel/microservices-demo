import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { getWebAutoInstrumentations } from '@opentelemetry/auto-instrumentations-web';
import { Resource } from '@opentelemetry/resources';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';

interface EndpointsConfig {
  apiArticlesEndpoint: string
  apiCartEndpoint: string
  otelEndpoint: string
}

let endpoints: EndpointsConfig = {
  apiArticlesEndpoint: "",
  apiCartEndpoint: "",
  otelEndpoint: ""
}


async function getEndpoints() : Promise<EndpointsConfig> {
  if (endpoints.otelEndpoint === "") {
    endpoints = await fetch("config/endpoints.json").then((response) => response.json())
  }
  return endpoints
}

export const provider = new WebTracerProvider({
  resource: Resource.default().merge(new Resource({
    // Replace with any string to identify this service in your system
    'service.name': 'front-user',
  })),
});

getEndpoints().then(endpointsConfig => {
  
  if(endpointsConfig.otelEndpoint !== "") {
    const traceExporter = new OTLPTraceExporter({
      url: endpointsConfig.otelEndpoint,
      headers: {
        'Content-Type': 'application/json',
        }    
    })
  
    provider.addSpanProcessor(new SimpleSpanProcessor(traceExporter));
  
    provider.register({
      contextManager: new ZoneContextManager(),
  
    });
  
    // Registering instrumentations
    registerInstrumentations({
      instrumentations: [
        getWebAutoInstrumentations(),
      ],
    });
    console.log('Tracing service started');
  }
})

