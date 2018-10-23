<!--
**목차**

API
- [Start API](#start-api)
- [On Calls API](#on-calls-api)
- [Action API](#action-api)

Resource
- [Elevator](#elevator)
- [Call](#call)
- [Command](#command)

---
-->

# API
HTTP로 통신한다.
- Server URL : http://api.welcome.kakao.com
- 초 당 40회가 넘는 API 호출은 서버가 응답하지 않을 수 있다.

<a name="start-api"></a>
## Start API
제어를 시작하기 위한 token을 발급하고 엘리베이터 상태 정보를 반환한다.

### Request

```http
POST /start/{user_key}/{problem_id}/{number_of_elevators}
```

| Parameter | Description |
| --------- | ----------- |
| `user_key` | 지원자가 사전에 지급 받은 user_key |
| `problem_id` | 제어를 시작할 문제 번호 (0 <= problem_id <= 2) |
| `number_of_elevators` | 사용할 엘리베이터 수 (1 <= count <= 4) |

### Response

| Status Code | 설명 |
| ----------- | ---- |
| 200 OK | 성공 |
| 400 Bad Request | `problem_id` 또는 `number_of_elevators`의 형식 또는 범위가 잘못됨 |
| 401 Unauthorized | `X-Auth-Token` Header가 잘못됨 |
| 403 Forbidden | `user_key`가 잘못되었거나 10초 이내에 생성한 토큰이 존재 |
| 500 Internal Server Error | 서버 에러, 문의 필요 |

성공일 때 response body는 JSON 형식이며 아래와 같다.

```json
{
  "token": "",
  "timestamp": 0,
  "elevators": [],
  "is_end": false
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | string | Token (Request header의 `X-Auth-Token` 값과 동일) |
| `timestamp` | integer | 현재 timestamp |
| `elevators` | array of [`Elevator`](#elevator) | 엘리베이터의 상태 |
| `is_end` | boolean | 모든 승객의 수송 완료 여부 |

#### Example

```json
{
  "token": "TVqpM5MX0amQqhrYKd3ZwMZn3Im6y4ovJwEa_A1-2d6mBpl4QhwJmmkrrvG4MsaD+mG44dL0aC0RNYL",
  "timestamp": 0,
  "elevators": [
    {
      "id": 0,
      "floor": 1,
      "passengers": [],
      "status": "STOPPED"
    },
    {
      "id": 1,
      "floor": 1,
      "passengers": [],
      "status": "STOPPED"
    }
  ],
  "is_end": false
}
```


<a name="on-calls-api"></a>
## On Calls API
현재 엘리베이터 상태와 현재 `timestamp` 기준으로 처리해야하는 call의 목록을 반환한다.

+ 승객이 탑승하면 call 목록에서 사라진다.
+ 탑승한 승객이 목적지가 아닌 층에서 하차하면, 하차한 층에서 목적지로 가는 새로운 `Call`이 call 목록에 추가된다.
+ 처리되지 않은 call은 call 목록에 계속 포함된다.

### Request

```http
GET /oncalls
X-Auth-Token: {Token}
```

| Name | Description |
| ---- | ----------- |
| `X-Auth-Token` | **Start API**의 결과로 받은 token |

### Response

| Status Code | 설명 |
| ----------- | ---- |
| 200 OK | 성공 |
| 401 Unauthorized | `X-Auth-Token` Header가 잘못됨 |
| 500 Internal Server Error | 서버 에러, 문의 필요 |

성공일 때 response body는 JSON 형식이며 아래와 같다.

```json
{
  "token": "",
  "timestamp": 0,
  "elevators": [],
  "calls": [],
  "is_end": false
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | string | Token (Request header의 `X-Auth-Token` 값과 동일)
| `timestamp` | integer | 현재 timestamp
| `elevators` | array of [`Elevator`](#elevator) | 엘리베이터의 상태
| `calls` | array of [`Call`](#call) | 현재 timestamp 기준으로 엘리베이터에 탑승하지 못한 call의 목록
| `is_end` | boolean | 모든 승객의 수송 완료 여부


#### Example

```json
{
  "token": "TVqpM5MX0amQqhrYKd3ZwMZn3Im6y4ovJwEa_A1-2d6mBpl4QhwJmmkrrvG4MsaD+mG44dL0aC0RNYL",
  "timestamp": 0,
  "elevators": [
    {
      "id": 0,
      "floor": 1,
      "passengers": [],
      "status": "STOPPED"
    },
    {
      "id": 1,
      "floor": 1,
      "passengers": [],
      "status": "STOPPED"
    }
  ],
  "calls": [
    {
      "id": 0,
      "timestamp": 0,
      "start": 6,
      "end": 1
    }
  ],
  "is_end": false
}
```


<a name="action-api"></a>
## Action API
엘리베이터에 명령을 실행한다.
+ `commands`에는 모든 엘리베이터의 명령이 포함되어 있어야 하며, 각 엘리베이터 당 하나의 명령만 전달할 수 있다.
+ 예를 들어 엘리베이터가 두 대인 경우 두 대 각각에 대한 명령이 있어야 한다.

### Request

```http
POST /action
X-Auth-Token: {Token}
Content-Type: application/json
```

| Name | Description |
| ---- | ----------- |
| `X-Auth-Token` | **Start API**의 결과로 받은 token |

```json
{
  "commands": []
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `commands` | array of [`Command`](#command) | 엘리베이터를 제어하기 위한 명령 |

#### Example

  ```json
  {
    "commands": [
      {
        "elevator_id": 0,
        "command": "ENTER",
        "call_ids": [0]
      },
      {
        "elevator_id": 1,
        "command": "STOP"
      }
    ]
  }
  ```

### Response

| Status Code | 설명 |
| ----------- | ---- |
| 200 OK | 성공 |
| 400 Bad Request | 해당 명령을 실행할 수 없음 (실행할 수 없는 상태일 때, 엘리베이터 수와 Command 수가 일치하지 않을 때, 엘리베이터 정원을 초과하여 태울 때) |
| 401 Unauthorized | `X-Auth-Token` Header가 잘못됨 |
| 500 Internal Server Error | 서버 에러, 문의 필요 |

성공일 때 response body는 JSON 형식이며 아래와 같다.

```json
{
  "token": "",
  "timestamp": 0,
  "elevators": [],
  "is_end": false
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `token` | string | Token (Request header의 `X-Auth-Token` 값과 동일)
| `timestamp` | integer | 현재 timestamp
| `elevators` | array of [`Elevator`](#elevator) | 엘리베이터의 상태
| `is_end` | boolean | 모든 승객의 수송 완료 여부

#### Example

```json
{
  "token": "TVqpM5MX0amQqhrYKd3ZwMZn3Im6y4ovJwEa_A1-2d6mBpl4QhwJmmkrrvG4MsaD+mG44dL0aC0RNYL",
  "timestamp": 8,
  "elevators": [
    {
      "id": 0,
      "floor": 6,
      "passengers": [
        {
          "id": 0,
          "timestamp": 0,
          "start": 6,
          "end": 1
        }
      ],
      "status": "OPENED"
    },
    {
      "id": 1,
      "floor": 1,
      "passengers": [],
      "status": "STOPPED"
    }
  ],
  "is_end": false
}
```
    

## Resource

<a name="elevator"></a>
### Elevator
엘리베이터 하나의 상태를 표현한다.

```json
{
  "id": 0,
  "floor": 1,
  "passengers": [],
  "status": "STOPPED"
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `id` | integer | 해당 엘리베이터 번호 |
| `floor` | integer | 해당 엘리베이터의 현재 위치(층) |
| `passengers` | array of [`Call`](#call) | 해당 엘리베이터에 타고 있는 승객들을 표현하는 `Call`의 목록 |
| `status` | string | 해당 엘리베이터 상태 |

<a name="call"></a>
### Call
`Call`을 표현한다.

```json
{
  "id": 0,
  "timestamp": 0,
  "start": 1,
  "end": 2
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `id` | integer | `Call`의 고유 번호 |
| `timestamp` | integer | 해당 `Call`이 발생한 `timestamp` |
| `start` | integer | 출발 층 |
| `end` | integer | 가려는 층 |

<a name="command"></a>
### Command
엘리베이터 제어 명령을 표현한다.

```json
{
  "elevator_id": 0,
  "command": "ENTER",
  "call_ids": [0]
}
```

| Name | Type | Description |
| ---- | ---- | ----------- |
| `elevator_id` | integer | 명령을 실행할 엘리베이터의 번호 ([`Elevator`](#elevator)의 `id`) |
| `command` | string | 실행할 명령 (`STOP`, `OPEN`, `ENTER`, `EXIT`, `CLOSE`, `UP`, `DOWN` 중 하나) |
| `call_ids` | array of integer | *(optional)* 태우거나 내려줄 [`Call`](#call)의 `id` |
