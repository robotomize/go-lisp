#include <stdio.h>
#include <chrono>

int fib(int n) {
    if (n <= 1) {
        return n;
    }

    return fib(n-1) + fib(n-2);
}

int main() {
    auto start = std::chrono::high_resolution_clock::now();
    printf("fib 37: %d\n", fib(40));
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    printf("Time taken: %lld milliseconds\n", duration.count());
}