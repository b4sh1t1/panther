/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import path from 'path';
import { loadDotEnvVars, getAppTemplateParams } from '../scripts/utils';

loadDotEnvVars(path.resolve(__dirname, '.env.test'));
/**
 * Mock the server-side EJS-injected AWS configuration.
 * See `web/public/index.ejs`
 */
const { PANTHER_CONFIG } = getAppTemplateParams();

const scriptTag = document.createElement('script');
scriptTag.id = '__PANTHER_CONFIG__';
scriptTag.type = 'application/json';
scriptTag.innerHTML = JSON.stringify(PANTHER_CONFIG);
document.body.appendChild(scriptTag);
