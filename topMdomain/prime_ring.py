import random
from sympy import isprime


class PrimeRing:
    """
    Prime Ring Zp* = {1, 2, ..., p}
    Initialize with a random prime:
    `pr = PrimeRing(31)`
    
    Usage:
    1.  for i, val in enumerate(pr): 
            pass
    2.  use `next(pr)` to generate indices.

    Note:
    Decrement rand_idx if indices need to start from 0
    """

    def __init__(self, p):
        self.prime = p
        self.factors = []

        n = p - 1
        i = 2
        while i*i <= n:
            if n % i == 0:
                self.factors.append(i)
                n //= i
            else:
                i += 1
        if n > 1:
            self.factors.append(n)
        
        self.generator = self.rand_generator(p)
        self.idx = 1

    @staticmethod
    def find_prime(r):
        while True:
            if isprime(r):
                return r
            r += 1

    def rand_generator(self, p):
        while True:
            g = random.randint(2, p-1)
            if all(pow(g, (p-1)//q, p) != 1 for q in self.factors):
                print(f"prime: {p}\tgenerator: {g}")
                return g
            
    def __iter__(self):
        return self
            
    def __next__(self):

        if self.idx == 0:
            raise StopIteration
        
        self.idx *= self.generator
        self.idx %= self.prime

        if self.idx == 1:
            self.idx = 0
            return 1
        
        return self.idx
