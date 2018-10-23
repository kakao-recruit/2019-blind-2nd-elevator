import random
import bisect


def gen(pr, count, min_floor, max_floor, seconds, mean, var, func=None):
    def incoming():
        return min_floor, random.randint(min_floor + 1, max_floor)

    def outgoing():
        return random.randint(min_floor + 1, max_floor), min_floor

    def interfloor():
        return random.sample(range(min_floor + 1, max_floor + 1), 2)

    if not func:
        func = [incoming, outgoing, interfloor]

    for i in range(1, len(pr)):
        pr[i] += pr[i - 1]

    ret = []
    userid = 0
    ts = 0

    while userid < count:
        idx = bisect.bisect_left(pr, random.random())
        f1, f2 = func[idx]()
        ret.append('%d,%d,%d,%d' % (ts, userid, f1, f2))
        ts += max(0, int(random.gauss(mean, var)))
        userid += 1

    return ret


def appeach_mansion():
    data = gen(pr=[0.33, 0.33, 0.34], count=6, min_floor=1, max_floor=6, seconds=6, mean=1.0, var=1.0)
    open('p0.in', 'w').write('\n'.join(data))


def jayg_building():
    data = gen(pr=[0.025, 0.025, 0.95], count=200, min_floor=1, max_floor=25, seconds=100, mean=2.0, var=5.0)
    open('p1.in', 'w').write('\n'.join(data))


def ryan_tower():
    min_floor = 1
    max_floor = 25
    kakao_lobby = 13

    kakao_floor = list(set(range(kakao_lobby, max_floor + 1)) - set([kakao_lobby]))
    other_floor = list(set(range(min_floor + 1, max_floor + 1)) - set(kakao_floor) - set([kakao_lobby]))

    def from_first_to_lobby():
        return min_floor, kakao_lobby

    def from_lobby_to_kakao():
        return kakao_lobby, random.choice(kakao_floor)

    def from_lobby_to_first():
        return kakao_lobby, min_floor

    def from_kakao_to_first():
        return random.choice(kakao_floor), min_floor

    def from_first_to_kakao():
        return min_floor, random.choice(kakao_floor)

    def interfloor_kakao():
        return random.sample(kakao_floor, 2)

    def incoming_other():
        return min_floor, random.choice(other_floor)

    def outgoing_other():
        return random.choice(other_floor), min_floor

    def interfloor_other():
        return random.sample(other_floor, 2)

    def kakao_guest():
        if random.random() <= 0.5:
            return from_first_to_lobby()
        return from_lobby_to_first()

    def kakao_employee():
        p = random.random()
        if p <= 1.0 / 4.0:
            return from_first_to_lobby()
        if p <= 2.0 / 4.0:
            return from_kakao_to_first()
        if p <= 4.0 / 4.0:
            return interfloor_kakao()

    def other_employee():
        p = random.random()
        if p <= 0.495:
            return incoming_other()
        if p <= 0.99:
            return outgoing_other()
        return interfloor_other()

    data = gen(pr=[0.44, 0.33, 0.23], count=500, min_floor=1, max_floor=25, seconds=150, mean=3.3, var=5.0, func=[kakao_guest, kakao_employee, other_employee])
    open('p2.in', 'w').write('\n'.join(data))


if __name__ == '__main__':
    appeach_mansion()
    jayg_building()
    ryan_tower()
