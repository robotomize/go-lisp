#include <stdio.h>
#include <chrono>

int fib(int n) {
    if (n <= 1) {
        return n;
    }

    return fib(n-1) + fib(n-2);
}
// Для максимальной эффективности на данном примере лучше всего подходит -O3 и отключение дополнительных флагов отладчика.
// Вот наиболее эффективная сборка с помощью g++ или clang++:
//
//    g++ -O3 -march=native -mtune=native -o fib main.cpp
//
// или на clang++:
//
//    clang++ -O3 -march=native -mtune=native -o fib main.cpp
//
// Опции:
// -O3             Максимальная оптимизация
// -march=native   Включает поддержку всех инструкций данного процессора
// -mtune=native   Тюнинг под конкретный CPU
// -o fib          Имя выходного файла
//
// Запустить исполнение: ./fib


int main() {
    auto start = std::chrono::high_resolution_clock::now();
    printf("fib 37: %d\n", fib(40));
    auto end = std::chrono::high_resolution_clock::now();
    auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(end - start);
    printf("Time taken: %lld milliseconds\n", duration.count());
}