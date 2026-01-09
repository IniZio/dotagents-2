---
globs: api/**
alwaysApply: false
---
# 🚀 FormX Sino API Development Guidelines

## 📋 Overview

This guide contains critical coding standards and patterns for the FormX Sino NestJS API server. The API follows a modular architecture with clear separation of concerns, using MikroORM for database access and BullMQ for async job processing.

## 🎯 **CRITICAL REQUIREMENTS**

**⚠️ MANDATORY**: Always run `make -C api format` after making code changes to ensure proper formatting.

**⚠️ MANDATORY**: Always use Makefile commands instead of raw npm scripts for consistency and proper workflow.

**⚠️ CRITICAL - MIGRATIONS**: NEVER create migration files manually. ALWAYS use:
```bash
make migration NAME=your_migration_name
```
This ensures proper formatting and follows the project's workflow. The Makefile automatically formats code after migration generation.

### Makefile Commands (MANDATORY)

**✅ DO:**
```bash
# Format code
make -C api format

# Run tests
make -C api test

# Lint code
make -C api lint

# Generate migrations (from project root)
make migration NAME=create_user_table
make auto-migration NAME=add_user_email

# Run migrations
make migrate-up
make migrate-down
```

**❌ DON'T:**
```bash
# Don't use raw npm scripts directly
npm run format
npm run test
npm run migration:create -- --name create_user_table
```

**Why Makefile?**
- Consistent interface across all developers
- Ensures proper workflow (e.g., formatting after migration generation)
- Centralized command definitions
- Works with docker-compose for containerized operations

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

### Module-Based Architecture (MANDATORY)

The API is organized by feature modules, with each module containing:
- **Controllers**: HTTP endpoints and request/response handling
- **Services**: Business logic and orchestration
- **DTOs**: Data Transfer Objects for validation
- **Processors**: BullMQ queue job processors (if applicable)
- **Module**: NestJS module definition

**⚠️ CRITICAL**: All **Entities** (database models) MUST be declared in the `database/entities/` module, NOT in feature modules. This ensures proper entity registration and avoids circular dependencies.

### Module Structure

```
api/src/modules/[feature]/
├── [feature].module.ts          # NestJS module definition
├── [feature].controller.ts       # HTTP endpoints
├── [feature].service.ts          # Business logic
├── dto/
│   ├── create-[feature].dto.ts   # Create DTOs
│   ├── update-[feature].dto.ts   # Update DTOs
│   └── [feature].response.dto.ts # Response DTOs
└── processors/                    # Queue processors (if applicable)
    └── [feature]-processor.ts
```

### Database Entities Structure

```
api/src/database/entities/
├── index.ts                # Entity exports (register all entities here)
├── user.ts                 # User entity
├── collection.ts           # Collection entity
├── upload.ts               # Upload entity
└── invoice.ts              # Invoice entity (when created)
```

### Core Modules Structure

```
api/src/
├── modules/                      # Feature modules
│   ├── collections/              # Collections module
│   ├── uploads/                  # Uploads (Documents) module
│   ├── invoices/                 # Invoices module
│   ├── formx/                    # FormX integration module
│   └── storage/                  # File storage module
├── queue/                        # Queue configuration
│   ├── queue.module.ts           # BullMQ module setup
│   └── processors/               # Queue processors
│       ├── formx-extraction.processor.ts
│       ├── formx-polling.processor.ts
│       ├── image-download.processor.ts
│       ├── document-status-update.processor.ts
│       └── collection-stats-update.processor.ts
├── database/                     # Database configuration
├── auth/                         # Authentication
├── config/                       # Configuration
└── main.ts                       # Application entry point
```

## 📦 Required Modules (FormX Integration)

Based on the FormX integration specification, the following modules must be created:

### 1. Collections Module
- **Purpose**: Manage collections of uploads/invoices
- **Entities**: `Collection`
- **Endpoints**: CRUD operations for collections
- **Services**: Collection business logic, stats calculation

### 2. Uploads Module
- **Purpose**: Handle file uploads and document processing
- **Entities**: `Upload` (renamed from Document)
- **Endpoints**: Upload files, list uploads, get upload details, retry failed uploads
- **Services**: File validation, storage, document status calculation
- **Processors**: Document status update processor

### 3. Invoices Module
- **Purpose**: Manage extracted invoices from FormX
- **Entities**: `Invoice`
- **Endpoints**: List invoices, get invoice details, update invoice (manual editing), approve invoice, retry failed extraction
- **Services**: Invoice creation from FormX results, status calculation, data mapping

### 4. FormX Module
- **Purpose**: FormX API integration and webhook handling
- **Services**: 
  - `FormXService`: FormX API client (upload, poll, get result, download image)
  - `FormXWebhookService`: Handle FormX webhook callbacks
- **Processors**:
  - `FormXExtractionProcessor`: Upload files to FormX
  - `FormXPollingProcessor`: Poll FormX job status
  - `ImageDownloadProcessor`: Download images from FormX
- **Controllers**: Webhook endpoint for FormX callbacks

### 5. Storage Module
- **Purpose**: File storage abstraction
- **Services**: `FileStorageService` (save, delete, get file streams)
- **Configuration**: Storage paths, cloud storage integration

## 🔌 Service Layer Pattern (MANDATORY)

### Service Structure

**✅ DO:**
```typescript
// api/src/modules/collections/collections.service.ts
import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@mikro-orm/nestjs';
import { EntityRepository, EntityManager } from '@mikro-orm/postgresql';
import { Collection } from '../../database/entities/collection.entity';
import { CreateCollectionDto } from './dto/create-collection.dto';
import { UpdateCollectionDto } from './dto/update-collection.dto';

@Injectable()
export class CollectionsService {
  constructor(
    @InjectRepository(Collection)
    private readonly collectionRepository: EntityRepository<Collection>,
    private readonly em: EntityManager,
  ) {}

  async create(createDto: CreateCollectionDto): Promise<Collection> {
    const collection = this.collectionRepository.create(createDto);
    await this.em.persistAndFlush(collection);
    return collection;
  }

  async findAll(): Promise<Collection[]> {
    return this.collectionRepository.findAll();
  }

  async findOne(id: string): Promise<Collection> {
    return this.collectionRepository.findOneOrFail(id);
  }

  async update(id: string, updateDto: UpdateCollectionDto): Promise<Collection> {
    const collection = await this.findOne(id);
    this.collectionRepository.assign(collection, updateDto);
    await this.em.flush();
    return collection;
  }

  async remove(id: string): Promise<void> {
    const collection = await this.findOne(id);
    await this.em.removeAndFlush(collection);
  }
}
```

**❌ DON'T:**
```typescript
// Don't put business logic in controllers
@Controller('collections')
export class CollectionsController {
  @Post()
  async create(@Body() dto: CreateCollectionDto) {
    // BAD - Business logic in controller
    const collection = new Collection();
    collection.name = dto.name;
    await this.em.persistAndFlush(collection);
    return collection;
  }
}
```

## 🎮 Controller Pattern (MANDATORY)

### Controller Structure

**✅ DO:**
```typescript
// api/src/modules/collections/collections.controller.ts
import { Controller, Get, Post, Body, Patch, Param, Delete, HttpCode, HttpStatus } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';
import { CollectionsService } from './collections.service';
import { CreateCollectionDto } from './dto/create-collection.dto';
import { UpdateCollectionDto } from './dto/update-collection.dto';
import { CollectionResponseDto } from './dto/collection-response.dto';

@ApiTags('collections')
@Controller('collections')
export class CollectionsController {
  constructor(private readonly collectionsService: CollectionsService) {}

  @Post()
  @ApiOperation({ summary: 'Create a new collection' })
  @ApiResponse({ status: 201, type: CollectionResponseDto })
  async create(@Body() createDto: CreateCollectionDto): Promise<CollectionResponseDto> {
    const collection = await this.collectionsService.create(createDto);
    return CollectionResponseDto.fromEntity(collection);
  }

  @Get()
  @ApiOperation({ summary: 'List all collections' })
  @ApiResponse({ status: 200, type: [CollectionResponseDto] })
  async findAll(): Promise<CollectionResponseDto[]> {
    const collections = await this.collectionsService.findAll();
    return collections.map(c => CollectionResponseDto.fromEntity(c));
  }

  @Get(':id')
  @ApiOperation({ summary: 'Get collection by ID' })
  @ApiResponse({ status: 200, type: CollectionResponseDto })
  async findOne(@Param('id') id: string): Promise<CollectionResponseDto> {
    const collection = await this.collectionsService.findOne(id);
    return CollectionResponseDto.fromEntity(collection);
  }

  @Patch(':id')
  @ApiOperation({ summary: 'Update collection' })
  @ApiResponse({ status: 200, type: CollectionResponseDto })
  async update(
    @Param('id') id: string,
    @Body() updateDto: UpdateCollectionDto,
  ): Promise<CollectionResponseDto> {
    const collection = await this.collectionsService.update(id, updateDto);
    return CollectionResponseDto.fromEntity(collection);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiOperation({ summary: 'Delete collection' })
  async remove(@Param('id') id: string): Promise<void> {
    await this.collectionsService.remove(id);
  }
}
```

## 📊 Entity Pattern (MANDATORY)

**⚠️ CRITICAL**: All entities MUST be placed in `api/src/database/entities/`, NOT in feature modules. This ensures:
- Proper entity registration in MikroORM
- Avoids circular dependency issues
- Centralized database schema management
- Easier migration generation

### Base Class Pattern (MANDATORY)

**⚠️ CRITICAL**: All entities MUST extend either `Base` or `AuditableBase` from `api/src/database/entities/base.ts`. NEVER manually define `id`, `createdAt`, `updatedAt`, or audit fields.

**Base Classes:**
- **`Base`**: For system/internal tables (e.g., `FormXJob`, `Blob`, `BlobAttachment`, `User`)
  - Provides: `id`, `createdAt`, `updatedAt`, `updatedBy`
- **`AuditableBase`**: For user-created tables (e.g., `Collection`, `Invoice`, `Upload`)
  - Extends `Base` and adds: `createdBy`, `createdByName`, `updatedBy` (overrides), `updatedByName`

**✅ DO:**
```typescript
// api/src/database/entities/collection.entity.ts
import { Entity, Property, OneToMany, Collection as MikroCollection } from '@mikro-orm/core';
import { AuditableBase } from './base';

@Entity({ tableName: 'collection' })
export class Collection extends AuditableBase {
  @Property({ type: 'varchar', length: 255 })
  name!: string;

  @Property({ type: 'text', nullable: true })
  description?: string;

  @Property({ type: 'date' })
  date!: Date;

  @Property({ type: 'jsonb', nullable: true })
  stats?: {
    queue: number;
    processing: number;
    ready_for_review: number;
    done: number;
    reviewRequired: number;
  };

  @OneToMany('Upload', 'collection')
  uploads = new MikroCollection<any>(this);
}
```

```typescript
// api/src/database/entities/formx-job.entity.ts
import { Entity, Property, Enum, Index } from '@mikro-orm/core';
import { Base } from './base';

@Entity({ tableName: 'formx_job' })
@Index({ properties: ['entityType', 'entityId'] })
export class FormXJob extends Base {
  @Enum(() => FormXJobEntityType)
  entityType!: FormXJobEntityType;

  @Property({ type: 'uuid' })
  entityId!: string;

  @Property({ type: 'text', nullable: true })
  formxJobId?: string;

  @Property({ type: 'timestamp', nullable: true })
  completedAt?: Date;
}
```

**❌ DON'T:**
```typescript
// BAD - Manually defining fields that Base/AuditableBase provides
@Entity({ tableName: 'collection' })
export class Collection {
  @PrimaryKey({ type: 'uuid', defaultRaw: 'gen_random_uuid()' })
  id: string = uuidv4();

  @Property({ type: 'timestamp', defaultRaw: 'CURRENT_TIMESTAMP' })
  createdAt: Date = new Date();

  @Property({ type: 'timestamp', defaultRaw: 'CURRENT_TIMESTAMP', onUpdate: () => new Date() })
  updatedAt: Date = new Date();
}
```

**Entity Classification:**
- **User-created entities** (extend `AuditableBase`): `Collection`, `Invoice`, `Upload`
- **System entities** (extend `Base`): `FormXJob`, `Blob`, `BlobAttachment`, `User`

### Entity Structure

**❌ DON'T:**
```typescript
// BAD - Don't put entities in feature modules
// api/src/modules/collections/entities/collection.entity.ts
@Entity({ tableName: 'collection' })
export class Collection {
  // ...
}
```

**✅ DO - Register entities:**
```typescript
// api/src/database/entities/index.ts
import type { EntityClass } from '@mikro-orm/core';
import { User } from './user';
import { Collection } from './collection.entity';
import { Upload } from './upload.entity';

export const entities: EntityClass<any>[] = [User, Collection, Upload];
```

## 📝 DTO Pattern (MANDATORY)

### DTO Structure

**✅ DO:**
```typescript
// api/src/modules/collections/dto/create-collection.dto.ts
import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsOptional, IsDateString, MaxLength } from 'class-validator';

export class CreateCollectionDto {
  @ApiProperty({ example: 'January 2025 Invoices' })
  @IsString()
  @MaxLength(255)
  name!: string;

  @ApiProperty({ example: 'Monthly invoice processing', required: false })
  @IsOptional()
  @IsString()
  description?: string;

  @ApiProperty({ example: '2025-01-27' })
  @IsDateString()
  date!: string;

  @ApiProperty({ example: 'user@example.com' })
  @IsString()
  creator!: string;
}
```

```typescript
// api/src/modules/collections/dto/collection-response.dto.ts
import { ApiProperty } from '@nestjs/swagger';
import { Collection } from '../entities/collection.entity';

export class CollectionResponseDto {
  @ApiProperty()
  id!: string;

  @ApiProperty()
  name!: string;

  @ApiProperty({ required: false })
  description?: string;

  @ApiProperty()
  date!: string;

  @ApiProperty()
  creator!: string;

  @ApiProperty({ required: false })
  stats?: {
    queue: number;
    processing: number;
    ready_for_review: number;
    done: number;
    reviewRequired: number;
  };

  @ApiProperty()
  createdAt!: Date;

  @ApiProperty()
  updatedAt!: Date;

  static fromEntity(entity: Collection): CollectionResponseDto {
    const dto = new CollectionResponseDto();
    dto.id = entity.id;
    dto.name = entity.name;
    dto.description = entity.description;
    dto.date = entity.date.toISOString().split('T')[0];
    dto.creator = entity.creator;
    dto.stats = entity.stats;
    dto.createdAt = entity.createdAt;
    dto.updatedAt = entity.updatedAt;
    return dto;
  }
}
```

## 🔄 Queue Processor Pattern (MANDATORY)

### Queue Processor Structure

**✅ DO:**
```typescript
// api/src/queue/processors/formx-extraction.processor.ts
import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import { Injectable, Logger } from '@nestjs/common';
import { FormXService } from '../../modules/formx/formx.service';
import { UploadsService } from '../../modules/uploads/uploads.service';

export interface FormXExtractionJobData {
  uploadId: string;
  filePath: string;
  collectionId: string;
  retryCount?: number;
}

@Processor('formx-extraction')
@Injectable()
export class FormXExtractionProcessor extends WorkerHost {
  private readonly logger = new Logger(FormXExtractionProcessor.name);

  constructor(
    private readonly formxService: FormXService,
    private readonly uploadsService: UploadsService,
  ) {
    super();
  }

  async process(job: Job<FormXExtractionJobData>): Promise<void> {
    const { uploadId, filePath, collectionId } = job.data;

    this.logger.log(`Processing FormX extraction for upload ${uploadId}`);

    try {
      const formxResponse = await this.formxService.uploadDocument(filePath);

      await this.uploadsService.updateFormXIds(uploadId, {
        formxJobId: formxResponse.job_id,
        formxRequestId: formxResponse.request_id,
        status: 'processing',
      });

      await this.enqueuePollingJob(uploadId, formxResponse.job_id, formxResponse.request_id);

      this.logger.log(`FormX extraction job enqueued for upload ${uploadId}`);
    } catch (error) {
      this.logger.error(`FormX extraction failed for upload ${uploadId}:`, error);
      await this.uploadsService.markError(uploadId, error.message);
      throw error;
    }
  }

  private async enqueuePollingJob(uploadId: string, formxJobId: string, formxRequestId: string): Promise<void> {
    // Implementation to enqueue polling job
  }
}
```

### FormX Polling Processor

**✅ DO:**
```typescript
// api/src/queue/processors/formx-polling.processor.ts
import { Processor, WorkerHost } from '@nestjs/bullmq';
import { Job } from 'bullmq';
import { Injectable, Logger } from '@nestjs/common';
import { FormXService } from '../../modules/formx/formx.service';
import { InvoicesService } from '../../modules/invoices/invoices.service';

export interface FormXPollingJobData {
  uploadId: string;
  formxJobId: string;
  formxRequestId: string;
  pollCount?: number;
  maxPolls?: number;
}

@Processor('formx-polling')
@Injectable()
export class FormXPollingProcessor extends WorkerHost {
  private readonly logger = new Logger(FormXPollingProcessor.name);
  private readonly MAX_POLLS = 60;
  private readonly POLL_INTERVAL = 5000;

  constructor(
    private readonly formxService: FormXService,
    private readonly invoicesService: InvoicesService,
  ) {
    super();
  }

  async process(job: Job<FormXPollingJobData>): Promise<void> {
    const { uploadId, formxJobId, formxRequestId, pollCount = 0, maxPolls = this.MAX_POLLS } = job.data;

    this.logger.log(`Polling FormX job ${formxJobId} (poll ${pollCount + 1}/${maxPolls})`);

    try {
      const status = await this.formxService.pollJobStatus(formxJobId);

      if (status.status === 'processing') {
        if (pollCount >= maxPolls) {
          throw new Error(`FormX extraction timeout after ${maxPolls} polls`);
        }

        await this.reEnqueuePolling(uploadId, formxJobId, formxRequestId, pollCount + 1);
        return;
      }

      if (status.status === 'completed') {
        const result = await this.formxService.getExtractionResult(formxRequestId);
        await this.invoicesService.createFromFormXResult(uploadId, result);
        await this.enqueueImageDownloads(uploadId, result);
        return;
      }

      if (status.status === 'failed') {
        throw new Error('FormX extraction failed');
      }
    } catch (error) {
      this.logger.error(`FormX polling failed for upload ${uploadId}:`, error);
      throw error;
    }
  }

  private async reEnqueuePolling(uploadId: string, formxJobId: string, formxRequestId: string, pollCount: number): Promise<void> {
    // Implementation to re-enqueue with delay
  }

  private async enqueueImageDownloads(uploadId: string, result: any): Promise<void> {
    // Implementation to enqueue image download jobs
  }
}
```

## 🔗 FormX Service Pattern (MANDATORY)

### FormX Service Structure

**✅ DO:**
```typescript
// api/src/modules/formx/formx.service.ts
import { Injectable, Logger } from '@nestjs/common';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom } from 'rxjs';
import { AppConfigService } from '../../config/config.service';

@Injectable()
export class FormXService {
  private readonly logger = new Logger(FormXService.name);
  private readonly baseUrl: string;
  private readonly token: string;
  private readonly extractorId: string;

  constructor(
    private readonly httpService: HttpService,
    private readonly configService: AppConfigService,
  ) {
    this.baseUrl = this.configService.get('formx.baseUrl');
    this.token = this.configService.get('formx.token');
    this.extractorId = this.configService.get('formx.extractorId');
  }

  async uploadDocument(filePath: string): Promise<{ job_id: string; request_id: string }> {
    const FormData = require('form-data');
    const fs = require('fs');

    const formData = new FormData();
    formData.append('image', fs.createReadStream(filePath));

    const response = await firstValueFrom(
      this.httpService.post(`${this.baseUrl}/v2/extract`, formData, {
        headers: {
          ...formData.getHeaders(),
          'X-WORKER-TOKEN': this.token,
          'X-WORKER-EXTRACTOR-ID': this.extractorId,
          'X-WORKER-ASYNC': 'true',
        },
      }),
    );

    return {
      job_id: response.data.job_id,
      request_id: response.data.request_id,
    };
  }

  async pollJobStatus(jobId: string): Promise<{ status: string }> {
    const response = await firstValueFrom(
      this.httpService.get(`${this.baseUrl}/v2/extract/jobs/${jobId}`, {
        headers: {
          'X-WORKER-TOKEN': this.token,
        },
      }),
    );

    return { status: response.data.status };
  }

  async getExtractionResult(requestId: string): Promise<any> {
    const response = await firstValueFrom(
      this.httpService.get(`${this.baseUrl}/v2/datalog/${requestId}/result`, {
        headers: {
          'X-WORKER-TOKEN': this.token,
        },
      }),
    );

    return response.data;
  }

  async downloadImage(requestId: string, pageNumber?: number, sliceNumber?: number): Promise<Buffer> {
    let url = `${this.baseUrl}/v2/datalog/${requestId}/file`;
    if (pageNumber !== undefined) {
      url += `?page=${pageNumber}`;
      if (sliceNumber !== undefined) {
        url += `&slice=${sliceNumber}`;
      }
    }

    const response = await firstValueFrom(
      this.httpService.get(url, {
        headers: {
          'X-WORKER-TOKEN': this.token,
        },
        responseType: 'arraybuffer',
      }),
    );

    return Buffer.from(response.data);
  }
}
```

## 🪝 Webhook Handler Pattern (MANDATORY)

### Webhook Controller Structure

**✅ DO:**
```typescript
// api/src/modules/formx/formx-webhook.controller.ts
import { Controller, Post, Body, Headers, HttpCode, HttpStatus, Logger } from '@nestjs/common';
import { ApiTags, ApiOperation } from '@nestjs/swagger';
import { FormXWebhookService } from './formx-webhook.service';

@ApiTags('formx')
@Controller('webhooks/formx')
export class FormXWebhookController {
  private readonly logger = new Logger(FormXWebhookController.name);

  constructor(private readonly webhookService: FormXWebhookService) {}

  @Post()
  @HttpCode(HttpStatus.OK)
  @ApiOperation({ summary: 'Handle FormX webhook callbacks' })
  async handleWebhook(
    @Body() payload: any,
    @Headers('x-formx-signature') signature?: string,
  ): Promise<{ received: boolean }> {
    this.logger.log('Received FormX webhook');

    try {
      await this.webhookService.verifySignature(payload, signature);
      await this.webhookService.processWebhook(payload);
      return { received: true };
    } catch (error) {
      this.logger.error('Webhook processing failed:', error);
      throw error;
    }
  }
}
```

```typescript
// api/src/modules/formx/formx-webhook.service.ts
import { Injectable, Logger } from '@nestjs/common';
import { InvoicesService } from '../invoices/invoices.service';
import { FormXService } from './formx.service';

@Injectable()
export class FormXWebhookService {
  private readonly logger = new Logger(FormXWebhookService.name);

  constructor(
    private readonly formxService: FormXService,
    private readonly invoicesService: InvoicesService,
  ) {}

  async verifySignature(payload: any, signature?: string): Promise<void> {
    // Implement signature verification if FormX provides webhook signing
    if (!signature) {
      this.logger.warn('Webhook signature missing');
    }
  }

  async processWebhook(payload: any): Promise<void> {
    const { event, data } = payload;

    switch (event) {
      case 'extraction.completed':
        await this.handleExtractionCompleted(data);
        break;
      case 'extraction.failed':
        await this.handleExtractionFailed(data);
        break;
      default:
        this.logger.warn(`Unknown webhook event: ${event}`);
    }
  }

  private async handleExtractionCompleted(data: any): Promise<void> {
    const { request_id, job_id } = data;
    this.logger.log(`Extraction completed for request ${request_id}`);

    const result = await this.formxService.getExtractionResult(request_id);
    // Process result and create invoices
  }

  private async handleExtractionFailed(data: any): Promise<void> {
    const { request_id, job_id, error } = data;
    this.logger.error(`Extraction failed for request ${request_id}: ${error}`);
    // Mark upload as failed
  }
}
```

## 📁 Module Registration Pattern (MANDATORY)

### Module Structure

**✅ DO:**
```typescript
// api/src/modules/collections/collections.module.ts
import { Module } from '@nestjs/common';
import { MikroOrmModule } from '@mikro-orm/nestjs';
import { CollectionsController } from './collections.controller';
import { CollectionsService } from './collections.service';
import { Collection } from '../../database/entities/collection.entity';

@Module({
  imports: [MikroOrmModule.forFeature([Collection])],
  controllers: [CollectionsController],
  providers: [CollectionsService],
  exports: [CollectionsService],
})
export class CollectionsModule {}
```

## 🚫 Common Pitfalls

### ❌ Don't Put Business Logic in Controllers

```typescript
// BAD
@Controller('collections')
export class CollectionsController {
  @Post()
  async create(@Body() dto: CreateCollectionDto) {
    const collection = new Collection();
    collection.name = dto.name;
    await this.em.persistAndFlush(collection);
    return collection;
  }
}

// GOOD
@Controller('collections')
export class CollectionsController {
  constructor(private readonly collectionsService: CollectionsService) {}

  @Post()
  async create(@Body() dto: CreateCollectionDto) {
    return this.collectionsService.create(dto);
  }
}
```

### ❌ Don't Forget Error Handling

```typescript
// BAD - No error handling
async process(job: Job) {
  const result = await this.formxService.uploadDocument(filePath);
  await this.updateDocument(result);
}

// GOOD - Proper error handling
async process(job: Job) {
  try {
    const result = await this.formxService.uploadDocument(filePath);
    await this.updateDocument(result);
  } catch (error) {
    this.logger.error(`Job failed:`, error);
    await this.markDocumentError(job.data.uploadId, error.message);
    throw error;
  }
}
```

### ❌ Don't Forget Query Invalidation After Mutations

```typescript
// BAD - No stats update after invoice creation
async createInvoice(data: CreateInvoiceDto) {
  const invoice = this.invoiceRepository.create(data);
  await this.invoiceRepository.persistAndFlush(invoice);
  return invoice;
}

// GOOD - Enqueue stats update job
async createInvoice(data: CreateInvoiceDto) {
  const invoice = this.invoiceRepository.create(data);
  await this.invoiceRepository.persistAndFlush(invoice);
  
  await this.documentStatusQueue.add('update-document-status', {
    documentId: invoice.uploadId,
  });
  
  return invoice;
}
```

### ❌ Don't Use Raw SQL When MikroORM Can Handle It

```typescript
// BAD - Raw SQL
const result = await this.em.execute('SELECT * FROM collections WHERE id = $1', [id]);

// GOOD - MikroORM query builder
const collection = await this.collectionRepository.findOneOrFail(id);
```

## 🔧 Development Workflow

### 1. Create New Module

```bash
# 1. Create module directory structure
mkdir -p api/src/modules/new-feature/dto

# 2. Create entity in database/entities (NOT in feature module)
touch api/src/database/entities/new-feature.entity.ts
# Register entity in api/src/database/entities/index.ts

# 3. Create DTOs
touch api/src/modules/new-feature/dto/create-new-feature.dto.ts
touch api/src/modules/new-feature/dto/update-new-feature.dto.ts
touch api/src/modules/new-feature/dto/new-feature-response.dto.ts

# 4. Create service
touch api/src/modules/new-feature/new-feature.service.ts

# 5. Create controller
touch api/src/modules/new-feature/new-feature.controller.ts

# 6. Create module
touch api/src/modules/new-feature/new-feature.module.ts

# 7. Register in app.module.ts
```

### 2. Code Review Checklist

- [ ] Uses module-based architecture
- [ ] Business logic in services, not controllers
- [ ] Proper DTOs with validation decorators
- [ ] Swagger/OpenAPI decorators on endpoints
- [ ] Error handling in services and processors
- [ ] Logging for important operations
- [ ] Queue processors handle errors gracefully
- [ ] Entities follow MikroORM patterns
- [ ] Module properly registered in app.module.ts
- [ ] DTOs have proper response mapping

### 3. Before Committing

```bash
# Format code
make -C api format

# Run tests
make -C api test

# Lint code
make -C api lint
```

## 📚 Key Dependencies

### Core Libraries
- **NestJS**: Framework
- **MikroORM**: Database ORM
- **BullMQ**: Job queue
- **class-validator**: DTO validation
- **class-transformer**: DTO transformation
- **@nestjs/swagger**: OpenAPI documentation

### Development Tools
- **TypeScript**: Type safety
- **ESLint**: Code linting
- **Prettier**: Code formatting

## 🧪 Testing Guidelines (MANDATORY)

### Test Coverage Requirements

**⚠️ CRITICAL**: All new code MUST have comprehensive unit tests with full coverage.

### Test File Structure

Place test files next to the code they test with `.spec.ts` suffix:

```
api/src/modules/[feature]/
├── [feature].service.ts
├── [feature].service.spec.ts
├── [feature].controller.ts
└── [feature].controller.spec.ts

api/src/queue/processors/
├── [feature]-processor.ts
└── [feature]-processor.spec.ts
```

### Testing Framework

- **Framework**: Vitest (configured in `package.json`)
- **Mocking**: Use `vi.fn()` and `vi.mock()` from Vitest
- **Coverage**: Run `make -C api test` (coverage reports available via npm scripts if needed)

### Test Pattern (MANDATORY)

**✅ DO:**
```typescript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { MyService } from './my.service';
import { EntityManager } from '@mikro-orm/postgresql';

// Mock external dependencies that have side effects
vi.mock('../external/service');
vi.mock('../../config/config.service');

describe('MyService', () => {
  let service: MyService;
  let mockDependency: any;

  beforeEach(() => {
    // Create fresh mocks for each test
    mockDependency = {
      method: vi.fn(),
    } as any;

    service = new MyService(mockDependency);
  });

  describe('methodName', () => {
    it('should handle success case', async () => {
      // Arrange
      mockDependency.method.mockResolvedValue({ data: 'result' });

      // Act
      const result = await service.methodName('input');

      // Assert
      expect(mockDependency.method).toHaveBeenCalledWith('input');
      expect(result).toEqual({ data: 'result' });
    });

    it('should handle error case', async () => {
      // Arrange
      mockDependency.method.mockRejectedValue(new Error('API error'));

      // Act & Assert
      await expect(service.methodName('input')).rejects.toThrow('API error');
    });

    it('should handle edge case (null/undefined)', async () => {
      mockDependency.method.mockResolvedValue(null);
      const result = await service.methodName('input');
      expect(result).toBeNull();
    });
  });
});
```

### What to Test

1. **Services**: All public methods with various scenarios
   - ✅ Success cases
   - ✅ Error cases
   - ✅ Edge cases (null, undefined, empty arrays)
   - ✅ Validation logic

2. **Processors**: Job processing logic
   - ✅ Successful job processing
   - ✅ Error handling and retries
   - ✅ Status transitions
   - ✅ Edge cases (missing data, invalid data)

3. **Controllers**: Request/response handling
   - ✅ Valid requests
   - ✅ Invalid requests (validation errors)
   - ✅ Authorization checks
   - ✅ Error responses

### Mocking Patterns

#### MikroORM EntityManager

```typescript
let em: EntityManager;
let mockFork: EntityManager;

beforeEach(() => {
  mockFork = {
    findOne: vi.fn(),
    findOneOrFail: vi.fn(),
    find: vi.fn(),
    create: vi.fn(),
    persistAndFlush: vi.fn(),
    flush: vi.fn(),
  } as any;

  em = {
    fork: vi.fn().mockReturnValue(mockFork),
  } as any;

  service = new MyService(em);
});
```

#### BullMQ Queue

```typescript
let queue: Queue;

beforeEach(() => {
  queue = {
    add: vi.fn().mockResolvedValue({ id: '1' }),
  } as any;

  processor = new MyProcessor(service, queue);
});
```

#### External Services

```typescript
// Use vi.mock() at the top of the file for services with side effects
vi.mock('../../modules/formx/formx.service');
vi.mock('../../config/config.service');

describe('MyProcessor', () => {
  let formxService: any;  // Use 'any' when mocked

  beforeEach(() => {
    formxService = {
      uploadDocument: vi.fn(),
      pollJobStatus: vi.fn(),
    } as any;

    processor = new MyProcessor(formxService);
  });
});
```

### Running Tests

```bash
# Run all tests (use Makefile)
make -C api test

# For advanced test options, check package.json scripts
# but prefer Makefile commands when available
```

### Test Checklist

Before committing code, ensure:

- [ ] All new services have unit tests
- [ ] All new processors have unit tests
- [ ] All new controllers have unit tests
- [ ] Tests cover success, error, and edge cases
- [ ] Tests use proper mocking (no real DB/API calls)
- [ ] Tests run successfully (`make -C api test`)
- [ ] Code coverage is maintained/improved

### Test Quality Standards

1. **Isolation**: Each test should be independent
2. **Clarity**: Test names should clearly describe what they test
3. **Arrange-Act-Assert**: Follow AAA pattern
4. **Mock Dependencies**: Don't make real API calls or DB queries
5. **Fast**: Tests should run quickly (no long delays or timeouts)

### Examples

See these files for reference implementations:
- `src/config/config.service.spec.ts` - Service testing
- `src/modules/uploads/uploads.service.spec.ts` - MikroORM service testing
- `src/queue/processors/formx-polling.processor.spec.ts` - Processor testing
- `src/queue/processors/formx-poll-scheduler.processor.spec.ts` - Scheduler testing

## 🎯 **REMINDER**

**⚠️ CRITICAL**: Always run `make -C api format` after making code changes.

**🧪 Testing**: Write comprehensive unit tests for all new code before committing.

**🏗️ Architecture**: Follow module-based architecture - separate controllers, services, entities, and DTOs.

**📦 Modules**: Organize code by feature modules, not by file type.

**🔄 Queues**: Use BullMQ processors for async operations. Handle errors gracefully.

**🔌 Services**: Business logic goes in services, not controllers.

**📝 DTOs**: Always use DTOs for request/response validation and transformation.

**🪝 Webhooks**: Implement webhook handlers for FormX callbacks (optional, polling is primary).

**📊 Entities**: Use MikroORM entities with proper relationships and indexes.
