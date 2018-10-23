import requests


url = 'http://localhost:8000'


def start(user, problem, count):
    uri = url + '/start' + '/' + user + '/' + str(problem) + '/' + str(count)
    return requests.post(uri).json()


def oncalls(token):
    uri = url + '/oncalls'
    return requests.get(uri, headers={'X-Auth-Token': token}).json()


def action(token, cmds):
    uri = url + '/action'
    return requests.post(uri, headers={'X-Auth-Token': token}, json={'commands': cmds}).json()


def p0_simulator():
    user = 'tester'
    problem = 0
    count = 2

    ret = start(user, problem, count)
    token = ret['token']
    print('Token for %s is %s' % (user, token))

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'UP'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'UP'}, {'elevator_id': 1, 'command': 'OPEN'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'UP'}, {'elevator_id': 1, 'command': 'ENTER', 'call_ids': [2, 3]}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'STOP'}, {'elevator_id': 1, 'command': 'CLOSE'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'OPEN'}, {'elevator_id': 1, 'command': 'UP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'ENTER', 'call_ids': [0, 1]}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'CLOSE'}, {'elevator_id': 1, 'command': 'OPEN'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'DOWN'}, {'elevator_id': 1, 'command': 'EXIT', 'call_ids': [2]}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'STOP'}, {'elevator_id': 1, 'command': 'CLOSE'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'OPEN'}, {'elevator_id': 1, 'command': 'UP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'EXIT', 'call_ids': [1]}, {'elevator_id': 1, 'command': 'UP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'CLOSE'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'DOWN'}, {'elevator_id': 1, 'command': 'OPEN'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'STOP'}, {'elevator_id': 1, 'command': 'EXIT', 'call_ids': [3]}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'OPEN'}, {'elevator_id': 1, 'command': 'CLOSE'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'ENTER', 'call_ids': [4]}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'CLOSE'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'DOWN'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'STOP'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'OPEN'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'EXIT', 'call_ids': [0, 4]}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'ENTER', 'call_ids': [5]}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'CLOSE'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'UP'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'STOP'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'OPEN'}, {'elevator_id': 1, 'command': 'STOP'}])

    oncalls(token)
    action(token, [{'elevator_id': 0, 'command': 'EXIT', 'call_ids': [5]}, {'elevator_id': 1, 'command': 'STOP'}])


if __name__ == '__main__':
    p0_simulator()
