---
globs: portal/**
alwaysApply: false
---
# 🚀 FormX Sino Portal Development Guidelines

## 📋 Overview

This guide contains critical coding standards and patterns for the FormX Sino Portal React application. The portal follows a modular, feature-based architecture with clear separation of concerns.

## 🎯 **CRITICAL REQUIREMENTS**

**⚠️ MANDATORY**: Always run `make format` in the portal directory after making code changes to ensure TypeScript correctness.

**⚠️ MANDATORY**: Always use Makefile commands instead of raw npm scripts for consistency and proper workflow.

### Makefile Commands (MANDATORY)

**✅ DO:**
```bash
# Format code
make -C portal format

# Run tests
make -C portal test

# Lint code
make -C portal lint

# Generate API types
make -C portal generate

# Run development server
make -C portal run

# Build for production
make -C portal build
```

**❌ DON'T:**
```bash
# Don't use raw npm scripts directly
npm run format:fix
npm run test
npm run generate-api-types
npm run dev
```

**Why Makefile?**
- Consistent interface across all developers
- Ensures proper workflow and dependencies
- Centralized command definitions
- Easy to discover available commands

## ⚙️ Configuration Pattern (MANDATORY)

### config.yaml Pattern

Both API and Portal use a unified `config.yaml` pattern with environment variable substitution via `envsubst` for local development and Docker deployments.

**⚠️ CRITICAL**: 
- Configuration MUST use `config.yaml` files (not environment variables directly) for better schema validation
- Configuration files MUST use `.template` suffix (e.g., `config.yaml.template`) with `${ENV_VAR}` syntax
- Local development uses `envsubst` to generate `config.yaml` from templates
- Config MUST be validated with Zod schemas for type safety
- Never commit generated `config.yaml` files (only templates)

**Why config.yaml?**
- Enables schema validation with Zod for type safety
- Centralized configuration management
- Environment variable substitution via `envsubst` for sensitive values
- Consistent pattern across API and Portal

## 🏗️ Architecture Patterns

### ViewModel Pattern (MANDATORY)

The portal follows a ViewModel pattern where business logic, data fetching, and state management are encapsulated in ViewModel hooks. Components receive ViewModel properties as spread props and contain only UI logic.

### Feature-Based Modular Architecture (MANDATORY)

The portal is organized by features, with each feature containing:
- **Components**: Pure UI components that receive ViewModel props
- **ViewModels**: Business logic, state management, and data orchestration
- **Queries**: React Query hooks for data fetching (used by ViewModels)
- **Services**: API client wrappers (used by Queries)
- **Actions**: Mutation actions (used by ViewModels)
- **Types**: TypeScript types and interfaces

### Feature Structure

```
portal/src/features/[feature]/
├── components/           # Feature-specific UI components
│   └── [Component].tsx
├── viewmodels/          # ViewModel hooks (business logic)
│   └── use[Feature]ViewModel.ts
├── queries/             # React Query hooks for data fetching
│   └── [feature]Queries.ts
├── services/            # API client wrappers
│   └── [feature]Service.ts
├── actions/             # Mutation actions
│   └── [feature]Actions.ts
├── types/               # TypeScript types
│   └── [feature]Types.ts
└── index.ts             # Public exports
```

### Core Structure

```
portal/src/core/
├── api/                 # API client configuration
│   ├── client.ts        # Base API client setup
│   └── types.ts         # OpenAPI generated types
├── components/          # Shared UI components (shadcn/ui)
├── lib/                 # Utility functions
│   ├── auth.ts          # Auth utilities
│   ├── config.ts         # Config utilities
│   └── utils.ts         # General utilities
└── providers/           # React context providers
    ├── auth/            # Auth provider
    └── config/           # Config provider
```

## 🔌 API Client Pattern (MANDATORY)

### OpenAPI React Query Integration

**✅ DO:**
```tsx
// portal/src/core/api/client.ts
import { useQueryClient } from "@tanstack/react-query";
import { useOpenApiQuery } from "openapi-react-query";
import { useConfig } from "@/components/providers/config";
import { useAuth } from "@/components/providers/auth";
import type { paths } from "./types";

export function useApiClient() {
  const { config } = useConfig();
  const { getAccessToken } = useAuth();

  return useOpenApiQuery<paths>({
    baseUrl: config.api.endpoint,
    getAuthHeaders: async () => {
      const token = await getAccessToken();
      return {
        Authorization: `Bearer ${token}`,
      };
    },
  });
}
```

### Service Layer Pattern

**✅ DO:**
```tsx
// portal/src/features/collections/services/collectionsService.ts
import { useApiClient } from "@/core/api/client";
import type { paths } from "@/core/api/types";

type GetCollectionsResponse = paths["/api/collections"]["get"]["responses"]["200"]["content"]["application/json"];

export function useCollectionsService() {
  const apiClient = useApiClient();

  const getCollections = () => {
    return apiClient.query({
      url: "/api/collections",
      method: "get",
    });
  };

  const getCollection = (id: string) => {
    return apiClient.query({
      url: "/api/collections/{id}",
      method: "get",
      params: { path: { id } },
    });
  };

  return {
    getCollections,
    getCollection,
  };
}
```

### Queries Pattern for Data Fetching

**✅ DO:**
```tsx
// portal/src/features/collections/queries/collectionsQueries.ts
import { useQuery } from "@tanstack/react-query";
import { useCollectionsService } from "../services/collectionsService";

export function useCollectionsQuery() {
  const { getCollections } = useCollectionsService();

  return useQuery({
    queryKey: ["collections"],
    queryFn: async () => {
      const response = await getCollections();
      return response.data?.collections ?? [];
    },
  });
}

export function useCollectionQuery(id: string) {
  const { getCollection } = useCollectionsService();

  return useQuery({
    queryKey: ["collections", id],
    queryFn: async () => {
      const response = await getCollection(id);
      return response.data;
    },
    enabled: !!id,
  });
}
```

### Actions Pattern for Mutations

**✅ DO:**
```tsx
// portal/src/features/collections/actions/collectionsActions.ts
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/core/api/client";
import { toast } from "sonner";

export function useCreateCollection() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: { name: string; description?: string; date: string }) => {
      const response = await apiClient.mutate({
        url: "/api/collections",
        method: "post",
        body: data,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
      toast.success("Collection created successfully");
    },
    onError: (error) => {
      toast.error("Failed to create collection");
      console.error("Create collection error:", error);
    },
  });
}

export function useUpdateCollection() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<Collection> }) => {
      const response = await apiClient.mutate({
        url: "/api/collections/{id}",
        method: "patch",
        params: { path: { id } },
        body: data,
      });
      return response.data;
    },
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
      queryClient.invalidateQueries({ queryKey: ["collections", variables.id] });
      toast.success("Collection updated successfully");
    },
    onError: (error) => {
      toast.error("Failed to update collection");
      console.error("Update collection error:", error);
    },
  });
}

export function useDeleteCollection() {
  const apiClient = useApiClient();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      const response = await apiClient.mutate({
        url: "/api/collections/{id}",
        method: "delete",
        params: { path: { id } },
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
      toast.success("Collection deleted successfully");
    },
    onError: (error) => {
      toast.error("Failed to delete collection");
      console.error("Delete collection error:", error);
    },
  });
}
```

### ViewModel Pattern

**✅ DO:**
```tsx
// portal/src/features/collections/viewmodels/useCollectionsViewModel.ts
import { useCallback } from "react";
import { useCollectionsQuery } from "../queries/collectionsQueries";
import { useCreateCollection, useUpdateCollection, useDeleteCollection } from "../actions/collectionsActions";
import { toast } from "sonner";

export interface CollectionsViewModel {
  collections: Collection[];
  isLoading: boolean;
  isError: boolean;
  handleCreate: (data: CreateCollectionInput) => void;
  handleUpdate: (id: string, data: UpdateCollectionInput) => void;
  handleDelete: (id: string) => void;
  handleRefresh: () => void;
}

export function useCollectionsViewModel(): CollectionsViewModel {
  const { data: collections = [], isLoading, isError, refetch } = useCollectionsQuery();
  const createCollection = useCreateCollection();
  const updateCollection = useUpdateCollection();
  const deleteCollection = useDeleteCollection();

  const handleCreate = useCallback((data: CreateCollectionInput) => {
    createCollection.mutate(data);
  }, [createCollection]);

  const handleUpdate = useCallback((id: string, data: UpdateCollectionInput) => {
    updateCollection.mutate({ id, data });
  }, [updateCollection]);

  const handleDelete = useCallback((id: string) => {
    deleteCollection.mutate(id);
  }, [deleteCollection]);

  const handleRefresh = useCallback(() => {
    refetch();
  }, [refetch]);

  return {
    collections,
    isLoading,
    isError,
    handleCreate,
    handleUpdate,
    handleDelete,
    handleRefresh,
  };
}
```

**Component Structure:**
```tsx
// portal/src/features/collections/components/CollectionsList.tsx
import { Card } from "@/components/ui/card";
import { CollectionsViewModel } from "../viewmodels/useCollectionsViewModel";

export interface CollectionsListProps extends CollectionsViewModel {}

export function CollectionsList({
  collections,
  isLoading,
  handleCreate,
  handleUpdate,
  handleDelete,
}: CollectionsListProps) {
  if (isLoading) {
    return <div className="p-4">Loading...</div>;
  }

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
      {collections.map((collection) => (
        <Card key={collection.id} className="p-6">
          <h3 className="text-lg font-semibold mb-2">{collection.name}</h3>
          <p className="text-sm text-muted-foreground">{collection.description}</p>
        </Card>
      ))}
    </div>
  );
}
```

**Route Integration:**
```tsx
// portal/src/routes/collections/index.tsx
import { createFileRoute } from "@tanstack/react-router";
import { CollectionsList } from "@/features/collections/components/CollectionsList";
import { useCollectionsViewModel } from "@/features/collections/viewmodels/useCollectionsViewModel";

export const Route = createFileRoute("/collections/")({
  component: CollectionsPage,
});

function CollectionsPage() {
  const viewModel = useCollectionsViewModel();

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-3xl font-bold mb-6">Collections</h1>
      <CollectionsList {...viewModel} />
    </div>
  );
}
```

**❌ DON'T:**
```tsx
// Don't put API calls or queries directly in components
export function CollectionsList() {
  const { data: collections, isLoading } = useCollectionsQuery();
  
  // ...
}

// Don't pass ViewModel as object prop
function CollectionsPage() {
  const viewModel = useCollectionsViewModel();
  return <CollectionsList viewModel={viewModel} />;
}
```

## 📁 File Organization

### Feature Module Export Pattern

**✅ DO:**
```tsx
// portal/src/features/collections/index.ts
export { CollectionsList } from "./components/CollectionsList";
export { CollectionCard } from "./components/CollectionCard";
export { useCollectionsViewModel } from "./viewmodels/useCollectionsViewModel";
export type { CollectionsViewModel } from "./viewmodels/useCollectionsViewModel";
export { useCollectionsQuery, useCollectionQuery } from "./queries/collectionsQueries";
export { useCreateCollection, useUpdateCollection, useDeleteCollection } from "./actions/collectionsActions";
export type { Collection } from "./types/collectionsTypes";
```

### Type Definitions

**✅ DO:**
```tsx
// portal/src/features/collections/types/collectionsTypes.ts
import type { paths } from "@/core/api/types";

type CollectionResponse = paths["/api/collections/{id}"]["get"]["responses"]["200"]["content"]["application/json"];

export type Collection = CollectionResponse;

export interface CreateCollectionInput {
  name: string;
  description?: string;
  date: string;
}

export interface UpdateCollectionInput {
  name?: string;
  description?: string;
  date?: string;
}
```

## 🔄 State Management

### React Query for Server State

- **Queries**: Use `useQuery` for fetching data
- **Mutations**: Use `useMutation` for create/update/delete operations
- **Query Invalidation**: Use `queryClient.invalidateQueries()` after mutations
- **Optimistic Updates**: Use `onMutate` for optimistic UI updates when appropriate

### Local State

- Use `useState` for component-local UI state (modals, form inputs, etc.)
- Use `useReducer` for complex local state logic
- Avoid global state management libraries unless necessary

## 🧩 Component Guidelines

### Use shadcn/ui Components

**✅ DO:**
```tsx
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";

<Button variant="default" size="lg">
  Submit
</Button>
```

**❌ DON'T:**
```tsx
// Don't create custom button components when shadcn/ui exists
<button className="custom-button">Click me</button>
```

### Styling with Tailwind CSS

**✅ DO:**
```tsx
<div className="flex items-center gap-4 p-6 bg-card border border-border rounded-lg">
  <h2 className="text-xl font-semibold text-foreground">Title</h2>
</div>
```

**❌ DON'T:**
```tsx
// Don't use inline styles or CSS modules for new components
<div style={{ padding: "24px", backgroundColor: "#fff" }}>
```

## 🚫 Common Pitfalls

### ❌ Don't Put API Calls in Components

```tsx
// BAD
const Component = () => {
  const [data, setData] = useState();
  
  useEffect(() => {
    fetch('/api/data').then(r => r.json()).then(setData);
  }, []);
}

// GOOD
const Component = ({ data }: { data: Data }) => {
  return <div>{data}</div>;
}
```

### ❌ Don't Mix Concerns

```tsx
// BAD - Component handles API calls, state, and UI
const Component = () => {
  const [data, setData] = useState();
  const [loading, setLoading] = useState(false);
  
  const handleSubmit = async () => {
    setLoading(true);
    const response = await fetch('/api/data', { method: 'POST', body: JSON.stringify(data) });
    const result = await response.json();
    setData(result);
    setLoading(false);
  };
  
  return <form onSubmit={handleSubmit}>...</form>;
}

// GOOD - ViewModel pattern with separation of concerns
const Component = ({ data, isLoading, handleSubmit }: ComponentViewModel) => {
  return <form onSubmit={handleSubmit}>...</form>;
}

function Page() {
  const viewModel = useComponentViewModel();
  return <Component {...viewModel} />;
}
```

### ❌ Don't Break ViewModel Pattern

```tsx
// BAD - Passing ViewModel object instead of spreading
function Page() {
  const viewModel = useCollectionsViewModel();
  return <CollectionsList viewModel={viewModel} />;
}

// BAD - Using queries directly in components
function CollectionsList() {
  const { data: collections } = useCollectionsQuery();
  return <div>{collections.map(...)}</div>;
}

// GOOD - Spreading ViewModel properties as props
function Page() {
  const viewModel = useCollectionsViewModel();
  return <CollectionsList {...viewModel} />;
}
```

### ❌ Don't Forget Query Invalidation

```tsx
// BAD - Mutation doesn't invalidate queries
const useCreateCollection = () => {
  return useMutation({
    mutationFn: createCollection,
    // Missing onSuccess with query invalidation
  });
}

// GOOD - Properly invalidates queries
const useCreateCollection = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: createCollection,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["collections"] });
    },
  });
}
```

### ❌ Don't Use localStorage for API Data

```tsx
// BAD - Using localStorage instead of API
const getCollections = () => {
  const data = localStorage.getItem('collections');
  return data ? JSON.parse(data) : [];
}

// GOOD - Using queries through ViewModel
const viewModel = useCollectionsViewModel();
const { collections } = viewModel;
```

### ❌ Don't Create Duplicate API Client Instances

```tsx
// BAD - Creating multiple API clients
const Component1 = () => {
  const client1 = createApiClient();
}

const Component2 = () => {
  const client2 = createApiClient();
}

// GOOD - Using shared hook
const Component1 = () => {
  const apiClient = useApiClient();
}

const Component2 = () => {
  const apiClient = useApiClient();
}
```

## 🔧 Development Workflow

### 1. Create New Feature

```bash
# 1. Create feature directory structure
mkdir -p portal/src/features/new-feature/{components,viewmodels,queries,services,actions,types}

# 2. Create service
touch portal/src/features/new-feature/services/newFeatureService.ts

# 3. Create queries
touch portal/src/features/new-feature/queries/newFeatureQueries.ts

# 4. Create actions
touch portal/src/features/new-feature/actions/newFeatureActions.ts

# 5. Create ViewModel
touch portal/src/features/new-feature/viewmodels/useNewFeatureViewModel.ts

# 6. Create types
touch portal/src/features/new-feature/types/newFeatureTypes.ts

# 7. Create components
touch portal/src/features/new-feature/components/NewFeatureComponent.tsx

# 8. Create index.ts
touch portal/src/features/new-feature/index.ts

# 9. Create route
touch portal/src/routes/new-feature/index.tsx
```

### 2. Code Review Checklist

- [ ] Uses ViewModel pattern correctly
- [ ] Components receive ViewModel props via spread (not as object)
- [ ] ViewModels encapsulate business logic and state
- [ ] Queries are in queries folder (not hooks)
- [ ] API calls go through service layer
- [ ] Data fetching uses queries pattern
- [ ] Mutations use actions pattern
- [ ] Components are pure UI (no direct API calls or queries)
- [ ] Proper TypeScript types defined
- [ ] Query invalidation after mutations
- [ ] Error handling with toast notifications
- [ ] Loading states handled
- [ ] Uses shadcn/ui components
- [ ] Uses Tailwind CSS for styling

### 3. Before Committing

```bash
# Format code
make -C portal format

# Run tests
make -C portal test

# Lint code
make -C portal lint
```

## 📚 Key Dependencies

### Core Libraries
- **React 19**: UI framework
- **TanStack Router**: File-based routing
- **TanStack React Query**: Server state management
- **openapi-react-query**: Type-safe API client
- **shadcn/ui**: Component library
- **Tailwind CSS**: Styling
- **TypeScript**: Type safety
- **Authgear**: Authentication

### Development Tools
- **Vite**: Build tool
- **ESLint**: Code linting
- **TypeScript**: Type checking

## 🧪 Testing & Browser Testing

### Browser Testing (MANDATORY)

**⚠️ CRITICAL**: Always test portal features in the browser after implementation to ensure the flow works correctly.

**Test Credentials**:
- Username: `formx-sino`
- Password: `formx-sino`

**Testing Workflow**:
1. Start the development server: `make -C portal run`
2. Open browser and navigate to the portal URL
3. Login with test credentials (`formx-sino` / `formx-sino`)
4. Test the implemented feature end-to-end
5. Verify UI interactions, API calls, and data flow
6. Check for errors in browser console
7. Verify loading states and error handling

**Why Browser Testing?**
- Ensures real-world user experience
- Catches integration issues between components
- Validates authentication flow
- Verifies API integration works correctly
- Confirms UI/UX is working as expected

## 🎯 **REMINDER**

**⚠️ CRITICAL**: Always run `make format` in the portal directory after making code changes.

**🏗️ Architecture**: Follow ViewModel pattern - components receive ViewModel properties as spread props, not as objects.

**📊 ViewModels**: Encapsulate business logic, state management, and data orchestration. Use queries and actions internally.

**🔍 Queries**: Use queries folder (not hooks) for React Query data fetching hooks.

**🔌 API**: Use service layer pattern with openapi-react-query for type-safe API calls.

**🎨 Components**: Use shadcn/ui components and Tailwind CSS for styling. Components are pure UI.

**🔄 State**: Use React Query for server state (via queries), useState for local UI state.

**📁 Organization**: Organize code by features, not by file type.
