import * as React from 'react';
import { createRoot } from 'react-dom/client';
import {
    MemoryRouter,
} from "react-router-dom";
import ProLayout from './layout/ProLayout'
import { ConfigDataContextProvider, configDataContext } from "./views/radar/fc2/configDataReducer";
const container = document.getElementById('app');

// Create a root.
const root = createRoot(container);

// Initial render: Render an element to the root.
root.render(
    <MemoryRouter initialEntries={[
        // "/v2x/V2xPanel"
        // "/radar/fc2/turret"
        "/radar/fc2/turrentTable"
        // "/radar/fc2/config2"
        // "/v2x/v2xPanel"
    ]}>
        <ConfigDataContextProvider>
            <ProLayout />
        </ConfigDataContextProvider>
    </MemoryRouter>
);
