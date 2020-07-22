# Internet VPC
resource "aws_vpc" "go-micro-main" {
  cidr_block = "10.0.0.0/16"
  instance_tenancy = "default"
  enable_dns_support = "true"
  enable_dns_hostnames = "true"
  enable_classiclink = "false"

  tags =  {
    Version = "latest"
    Name = var.SERVICE_NAME
  }
}

# Subnets
resource "aws_subnet" "go-micro-main-public-1" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = "true"
  availability_zone = "ap-northeast-1a"

  tags = {
    Name = "${var.SERVICE_NAME}-public-1"
  }
}

resource "aws_subnet" "go-micro-main-public-2" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.2.0/24"
  map_public_ip_on_launch = "true"
  availability_zone = "ap-northeast-1c"

  tags = {
    Name = "${var.SERVICE_NAME}-public-2"
  }
}

resource "aws_subnet" "go-micro-main-public-3" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.3.0/24"
  map_public_ip_on_launch = "true"
  availability_zone = "ap-northeast-1d"

  tags = {
    Name = "${var.SERVICE_NAME}-public-3"
  }
}

resource "aws_subnet" "go-micro-main-private-1" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.4.0/24"
  map_public_ip_on_launch = "false"
  availability_zone = "ap-northeast-1a"

  tags = {
    Name = "${var.SERVICE_NAME}-private-1"
  }
}

resource "aws_subnet" "go-micro-main-private-2" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.5.0/24"
  map_public_ip_on_launch = "false"
  availability_zone = "ap-northeast-1c"

  tags = {
    Name = "${var.SERVICE_NAME}-private-2"
  }
}

resource "aws_subnet" "go-micro-main-private-3" {
  vpc_id = aws_vpc.go-micro-main.id
  cidr_block = "10.0.6.0/24"
  map_public_ip_on_launch = "false"
  availability_zone = "ap-northeast-1d"

  tags = {
    Name = "${var.SERVICE_NAME}-private-3"
  }
}

# IGW
resource "aws_internet_gateway" "go-micro-main-gw" {
  vpc_id = aws_vpc.go-micro-main.id
  
  tags = {
    Name = "${var.SERVICE_NAME}-gw"
  }
}

# route tables
resource "aws_route_table" "go-micro-main-public" {
  vpc_id = aws_vpc.go-micro-main.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.go-micro-main-gw.id
  }

  tags = {
    Name = "${var.SERVICE_NAME}-public-1"
  }
}

# route association public
resource "aws_route_table_association" "main-public-1-a" {
  subnet_id = aws_subnet.go-micro-main-public-1.id
  route_table_id = aws_route_table.go-micro-main-public.id
}
resource "aws_route_table_association" "main-public-2-a" {
  subnet_id = aws_subnet.go-micro-main-public-2.id
  route_table_id = aws_route_table.go-micro-main-public.id
}
resource "aws_route_table_association" "main-public-3-a" {
  subnet_id = aws_subnet.go-micro-main-public-3.id
  route_table_id = aws_route_table.go-micro-main-public.id
}
