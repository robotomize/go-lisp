; пока i больше нуля, a и b сдвигаются вместе
(define (fib-do n)
  (do ((i n (- i 1))
       (a 0 b)
       (b 1 (+ a b)))
      ((> i 0) a)))

(display "fib-do 1000")
(display (fib-do 1000))
