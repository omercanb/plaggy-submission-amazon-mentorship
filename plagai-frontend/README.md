# PlagAI Portal Frontend

Repo for the instructor portal for PlagAI, also know as Plaggy. For the live deployment head over to [plaggy.xyz](plaggy.xyz).

## Tech Used

- Next.js: Provides easy routing and also server-side rendering for future scalability
- MaterialUI: Easy to use and professional looking styled components
- TypeScript: Essential for a project of this size since without type management development would get out of hand quickly

## Features

- Assign homeworks which students can later access through the CLI
- View flagged submissions
- Manage sections and homeworks from a single dashboard

### Local Setup

- Install node.js v24.4.1 (or latest if not available) from [nodejs.org](https://nodejs.org/tr/download) or using nvm
- Clone the repository
- Run `npm install .`
- Make sure the backend is running and if necessary change the api url in `.env` or `/src/app/lib/api.ts`
- Run `npm run dev`
- Head over to `localhost:3000` and you should see the page

### Deployment

- Follow the same steps as local setup until `npm run dev`
- Instead, run `npm run build`. This will optimize the project and build it into static files.
- Once the build is done, run `npm run start` to serve both static and server-side rendered pages.
- Make sure the frontend is able to communicate with the backend. If not, the problem most likely is caused by the api url. Make sure it is /api/v1 since we are going to be using nginx as a reverse proxy.
- I recommend using pm2 to manage the two processes. To install, run `npm install -g pm2` on the machine.
- Running `pm2 start "npm run start -- -p 3000 -H 0.0.0.0" --name frontend` should start the frontend server.
