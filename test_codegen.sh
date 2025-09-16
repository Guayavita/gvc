#!/bin/bash

# Test script for Guayavita code generation functionality
# Tests binary compilation, JIT execution, and LLVM IR emission

set -e

echo "Building Guayavita compiler..."
go build -o ./bin/guayavita ./main.go

echo ""
echo "Testing different compilation modes with hello-simple.gvt..."

TEST_FILE="test-data/hello-simple.gvt"
OUTPUT_DIR="./bin"

# Ensure output directory exists
mkdir -p $OUTPUT_DIR

echo ""
echo "1. Testing syntax-only mode..."
./bin/guayavita compile --syntax-only $TEST_FILE

echo ""
echo "2. Testing LLVM IR emission..."
./bin/guayavita compile --emit-llvm -o $OUTPUT_DIR $TEST_FILE

echo ""
echo "3. Testing binary compilation..."
./bin/guayavita compile -o $OUTPUT_DIR $TEST_FILE

echo ""
echo "4. Testing JIT execution..."
./bin/guayavita compile --jit $TEST_FILE

echo ""
echo "5. Testing with custom target triple..."
./bin/guayavita compile --target x86_64-pc-linux-gnu --emit-llvm -o $OUTPUT_DIR $TEST_FILE

echo ""
echo "Checking generated files..."
if [ -f "$OUTPUT_DIR/hello-simple.ll" ]; then
    echo "✓ LLVM IR file generated: $OUTPUT_DIR/hello-simple.ll"
    echo "First few lines of LLVM IR:"
    head -n 10 "$OUTPUT_DIR/hello-simple.ll" || echo "Could not read LLVM IR file"
else
    echo "✗ LLVM IR file not found"
fi

echo ""
echo "All tests completed!"
echo ""
echo "Usage examples:"
echo "  Compile to binary:     ./bin/guayavita compile test-data/hello-simple.gvt"
echo "  Emit LLVM IR:         ./bin/guayavita compile --emit-llvm test-data/hello-simple.gvt"
echo "  JIT execution:        ./bin/guayavita compile --jit test-data/hello-simple.gvt"
echo "  Custom target:        ./bin/guayavita compile --target x86_64-pc-linux-gnu test-data/hello-simple.gvt"
echo "  Syntax only:          ./bin/guayavita compile --syntax-only test-data/hello-simple.gvt"