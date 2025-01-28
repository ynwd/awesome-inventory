# Awesome Inventory

Simple multi-tenant inventory management system.

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker

## Setup

1. Clone the repository
```bash
git clone https://github.com/ynwd/awesome-inventory.git
cd awesome-inventory
```

2. Initialize the database
```bash
docker compose up
```

3. Run the application
```bash
go run cmd/main.go
```

## Database Schema

The system uses a multi-tenant architecture with separated schemas for each tenant.

### System Tables (Main Database)

#### Tenants Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| DatabaseName | varchar(100) | Tenant database name |
| DatabaseHost | varchar(255) | Database host address |
| DatabasePort | integer | Database port number |
| DatabaseUser | varchar(100) | Database username |
| DatabasePass | varchar(255) | Database password (encrypted) |
| Status | varchar(20) | Tenant status |

#### Tenant Owners Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| Email | varchar(255) | Unique email address |
| Name | varchar(100) | Owner's name |
| Password | varchar(255) | Hashed password |
| CreatedAt | timestamp | Creation timestamp |
| UpdatedAt | timestamp | Last update timestamp |

### Tenant-Specific Tables (Per Tenant Database)

#### Users Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| Email | varchar(255) | Unique email address |
| Password | varchar(255) | Hashed password |
| FirstName | varchar(100) | User's first name |
| LastName | varchar(100) | User's last name |
| Role | varchar(50) | User role |
| IsActive | boolean | Account status |

#### Categories Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| Name | varchar(100) | Category name |
| Description | text | Category description |

#### Products Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| Name | varchar(255) | Product name |
| SKU | varchar(50) | Stock keeping unit |
| SupplierID | uuid | Foreign key to suppliers |
| Description | text | Product description |
| UnitPrice | decimal | Price per unit |

#### Suppliers Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| Name | varchar(255) | Supplier name |
| Code | varchar(50) | Supplier code |
| Contact | varchar(100) | Contact person |
| Email | varchar(255) | Contact email |
| Phone | varchar(50) | Contact phone |
| IsActive | boolean | Supplier status |

#### Inventory Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| ProductID | uuid | Foreign key to products |
| Quantity | integer | Current stock quantity |
| Location | varchar(100) | Storage location |

#### Inventory Transfers Table
| Column | Type | Description |
|:-------|:-----|:------------|
| ID | uuid | Primary key |
| InventoryID | uuid | Foreign key to inventory |
| Type | varchar(50) | Transfer type |
| Quantity | integer | Transfer quantity |
| SourceLocation | varchar(100) | Source location |
| DestLocation | varchar(100) | Destination location |
| TransferDate | timestamp | Transfer timestamp |
| Notes | text | Transfer notes |

## Testing

```bash
go test ./...
```

## License

MIT License