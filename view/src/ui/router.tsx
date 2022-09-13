import * as React from 'react';
import {
    Routes,
    Route,
    useLocation
} from "react-router-dom";

import Fc2Config2 from './views/radar/fc2/config2'
import Fc2Monitor from './views/radar/fc2/monitor'
import Fc2Turret from './views/radar/fc2/turret'
import Fc2DarkroomModule from './views/radar/fc2/darkroomModule'
import TurretTable from './views/radar/fc2/turrentTable'
import V2xPanel from './views/v2x/v2xPanel'

const _call = (_config) => {
    console.log("🚀 ~ file: router.tsx ~ line 18 ~ _config", _config)

}

export const AppRouter = () => (
    <Routes>
        <Route path="/" element={<Fc2Monitor />} />
        <Route path="radar">
            <Route path="fc2">
                <Route path="monitor" element={<Fc2Monitor />} />
                <Route path="turret" element={<Fc2Turret />} />
                <Route path="config2" element={<Fc2Config2 />} />
                <Route path="darkroom/module" element={<Fc2DarkroomModule />} />
                <Route path="turrentTable" element={<TurretTable />} />
            </Route>
        </Route>
        <Route path="v2x">
            <Route path="V2xPanel" element={<V2xPanel />} />
        </Route>
        <Route
            path="*"
            element={
                <main style={{ padding: "1rem" }}>
                    <p>There's nothing here!</p>
                    <p>{useLocation().pathname}</p>
                </main>
            }
        />
    </Routes >
);

