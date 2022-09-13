import turretApi from '../api/fc2/turretApi'
import numeral from 'numeral';

const moduleType = 1;
class TurretService {

    async getStatus() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: 'RED_angle'
        })
        return response;
    }

    async getStatusAck() {
        const response = await turretApi.commandAcknowledge(moduleType)  
        return parseStatus(response.commandPayload);
    }

    async setLaser(checked) {
        const response = await turretApi.command(
            checked ?
                {
                    moduleType: moduleType,
                    commandPayload: 'SET_laser'
                }
                :
                {
                    moduleType: moduleType,
                    commandPayload: 'RST_laser'
                })
        return response;
    }

    async start() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: 'SET_start'
        })
        return response;
    }

    async setAngle(values) {
        // console.log("🚀 ~ file: turretService.ts ~ line 44 ~ TurretService ~ setAngle ~ values", values)
        const angleFomat = '+000.00'
        
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: `SET_angle_${numeral(values.Polar.range).format(angleFomat)}_${numeral(values.EL.range).format(angleFomat)}_${numeral(values.AZ.range).format(angleFomat)}`
        })
        return response;
    }

    async setPartAngle(data) {
        const angleFomat = '+000.00'
        let cmd;
        switch (data.name) {
            case 'Polar':
                cmd = `SET_11111_${numeral(data.value).format(angleFomat)}`
                break;
            case 'EL':
                cmd = `SET_22222_${numeral(data.value).format(angleFomat)}`
                break;
            case 'AZ':
                cmd = `SET_33333_${numeral(data.value).format(angleFomat)}`
                break;
            default:
                throw (`unkown data.type: ${data.type}`)
                break;
        }
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: cmd
        })
        return response;
    }


    async setSpeed(values) {
        const speedFomat = '000.00'
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: `SET_speed_${numeral(values.Polar.maxSpeed).format(speedFomat)}_${numeral(values.EL.maxSpeed).format(speedFomat)}_${numeral(values.AZ.maxSpeed).format(speedFomat)}`
        })
        return response;
    }

    async stop() {
        const response = await turretApi.command({
            moduleType: moduleType,
            commandPayload: "SET_cease"
        })
        return response;
    }
}

function parseStatus(message) {
    const data = {
        'workStatus': 'idle',
        'data': {
            "laser": true,
            "Polar": {
                "angle": 1,
                "speed": 2
            },
            "EL": {
                "angle": 3,
                "speed": 4
            },
            "AZ": {
                "angle": 5,
                "speed": 6
            }
        }
    };

    data['data']['Polar']['angle'] = Number(message.substring(10, 17));
    data['data']['EL']['angle'] = Number(message.substring(18, 25));
    data['data']['AZ']['angle'] = Number(message.substring(26, 33));
    if ('RED' === message.substring(0, 3)) {
        data['workStatus'] = 'work';
    } else {
        data['workStatus'] = 'idle';
    }
    return data;
}
const turretService = new TurretService()
export default turretService