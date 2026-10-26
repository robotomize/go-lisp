; дерево вызовов: fib(n) = fib(n-1) + fib(n-2)
(define (fib n)
  (if (< n 2)
      n
      (+ (fib (- n 1))
         (fib (- n 2)))))

;(display "fib 30")
(display (fib 30))