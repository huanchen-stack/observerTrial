import csv

from prime_ring import PrimeRing

fin = "majestic_million.csv" 
fout = "urls.txt"

def load_urls(amount, csv_path, url_col):
    all_urls = []
    with open(csv_path, 'r') as csv_file:
        csv_reader = csv.reader(csv_file)

        # skip the header and most popular websites
        for i in range(1001):
            next(csv_reader)
        
        for i, row in enumerate(csv_reader):
            if i == 500000: break
            all_urls.append(row[url_col])
    
    # rand idx urls
    urls = []
    p = PrimeRing.find_prime(len(all_urls))
    pr = PrimeRing(p)  # generator for rand and evenly distributed indices
    with open(fout, 'w') as f_batch_input:
        for i, rand_idx in enumerate(pr):
            if len(urls) == amount: 
                break

            if not rand_idx-1 < len(all_urls):
                continue

            url = all_urls[rand_idx-1]
            print(rand_idx, url)
            f_batch_input.write(f"{url}\n")
            urls.append(url)


load_urls(10000, fin, url_col=2)