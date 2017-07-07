import { Injectable } from '@angular/core';
import { Project } from './project';
import { ProjectMember } from './project-member';

const PROJECTS = [
  {
    id: 1,
    projectName: 'Project01',
    creator: 'Tester',
    creationTime: new Date('2017-07-07 10:20:35'),
    isPublic: true,
    comment: 'something to say...'
  },
  {
    id: 2,
    projectName: 'Project02',
    creator: 'User1',
    creationTime: new Date('2017-07-07 09:20:35'),
    isPublic: false,
    comment: 'something to say...'
  },
  {
    id: 3,
    projectName: 'Project03',
    creator: 'Tester',
    creationTime: new Date('2017-07-07 11:33:45'),
    isPublic: true,
    comment: 'something to say...'
  }
];

const PROJECT_MEMBERS = [
  {
    projectId: 1,
    members: [
      {
        projectId: 1,
        userId: 1,
        username: 'Zhao'
      },
      {
        projectId: 1,
        userId: 3,
        username: 'Qian'
      },
      {
        projectId: 1,
        userId: 5,
        username: 'Sun'
      }
    ]
  },
  {
    projectId: 2,
    members: [
      {
        projectId: 1,
        userId: 3,
        username: 'Qian'
      }
    ]
  },
  {
    projectId: 3,
    members: [
      {
        projectId: 1,
        userId: 1,
        username: 'Zhao'
      },
      {
        projectId: 1,
        userId: 5,
        username: 'Sun'
      }
    ]
  }
];

@Injectable()
export class ProjectService {
  getProjects(): Promise<Project[]> {
    return Promise.resolve(PROJECTS);
  }

  getProjectMembers(): Promise<ProjectMember[]> {
    return Promise.resolve(PROJECT_MEMBERS);
  }
}