(define-macro (and . args)
  (if (null? args)
      #t
      (if (null? (cdr args))
          (car args)
          (list 'if (car args) (cons 'and (cdr args)) #f))))

(define-macro (or . args)
  (if (null? args)
      #f
      (list 'let
            (list (list 'value (car args)))
            (list 'if 'value 'value (cons 'or (cdr args))))))

; пока i больше нуля, a и b сдвигаются вместе
(define (fib-do n)
  (do ((i n (- i 1))
       (a 0 b)
       (b 1 (+ a b)))
      ((> i 0) a)))

(define (length xs)
  (if (null? xs)
      0
      (+ 1 (length (cdr xs)))))

(define (append xs ys)
  (if (null? xs)
      ys
      (cons (car xs) (append (cdr xs) ys))))

(define (reverse xs)
  (define (rev xs acc)
    (if (null? xs)
        acc
        (rev (cdr xs) (cons (car xs) acc))))
  (rev xs ()))

(define (map f xs)
  (if (null? xs)
      ()
      (cons (f (car xs)) (map f (cdr xs)))))

(define (filter pred xs)
  (if (null? xs)
      ()
      (if (pred (car xs))
          (cons (car xs) (filter pred (cdr xs)))
          (filter pred (cdr xs)))))

(define (accumulate op init xs)
  (if (null? xs)
      init
      (op (car xs) (accumulate op init (cdr xs)))))

(define (caar x)
  (car (car x)))

(define (cadr x)
  (car (cdr x)))

(define (len xs)
  (length xs))

(define (range a b)
  (if (> a b)
      ()
      (cons a (range (+ a 1) b))))

(define (foldl op acc xs)
  (if (null? xs)
      acc
      (foldl op (op acc (car xs)) (cdr xs))))

(define (foldr op init xs)
  (accumulate op init xs))

(define (reduce op xs)
  (if (null? xs)
      ()
      (foldl op (car xs) (cdr xs))))

(define (take n xs)
  (if (or (<= n 0) (null? xs))
      ()
      (cons (car xs) (take (- n 1) (cdr xs)))))

(define (interleave a b)
  (if (null? a) b
  (if (null? b) a
    (cons (car a)
      (cons (car b) (interleave (cdr a) (cdr b)))))))


(define (drop n xs)
  (if (or (<= n 0) (null? xs))
      xs
      (drop (- n 1) (cdr xs))))

(define (take-while pred xs)
  (if (or (null? xs) (not (pred (car xs))))
      ()
      (cons (car xs) (take-while pred (cdr xs)))))

(define (drop-while pred xs)
  (if (or (null? xs) (not (pred (car xs))))
      xs
      (drop-while pred (cdr xs))))

(define-macro (when test body)
  (list 'if test body))

(define-macro (unless test body)
  (list 'if test () body))

(define-macro (defn name params body)
  (list 'define (cons name params) body))

(define-macro (defmacro name params body)
  (list 'define-macro (cons name params) body))
