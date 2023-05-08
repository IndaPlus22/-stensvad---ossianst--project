# Specification
A specification for a 3D renderer and planet generator using OpenGL and the Go programming language.  

## Table of contents
* [Relevant Links](#relevant-links)
* [Branches](#branches)
* [Issues](#issues)
* [Commits](#commits)
* [Project](#project)
  * [Idea](#idea)
  * [Requrements](#requirements)
     * [Must have](#must-have)
     * [Should have](#should-have)
     * [Could have](#could-have)
     * [Will not have](#will-not-have)
  * [Feasibillity](#feasibillity)
  * [Areas of responsibility](#areas-of-responsibility)
  
## Relevant Links
* [Repo page](https://github.com/IndaPlus22/stensvad-ossianst-melvinbe-project)
* [Issues page](https://github.com/IndaPlus22/stensvad-ossianst-melvinbe-project/issues)
* [Project board](https://github.com/orgs/IndaPlus22/projects/1/views/1)

## Branches
* main - The main branch must be functional after every merge and remain protected from direct pushes. 

* dev - A branch for sub-branches to merge into. The dev branch will merge into the main branch when functional. Direct pushes to the dev branch are allowed.

* issue - An issue branch is created for each issue and merged into dev when solved. These branches will be named after their corresponding issues. The issue `#3 Add-example` would be fixed in branch `issue/3-add-example`. 

## Issues
Issues will be created for every feature to be added or change to make to the program. 

Small enough changes that do not reqire issues can be commited direcly to dev. 

Issues will have concise names describing the fix or feature, and a short description detailing what that would entail. 

The title should be written in future tense. 

Assignees will be assigned to every issue.

Issues with deadlines will be added under milestones. 

Issues should be fixed automatically by writing `[Ff]ix #3` in the commit message.

## Commits
Commit messages should be prefixed with a category of what change has been made and be written in future tense, like this: `feature: add example` or `bug fix: fix example`. 

Commits should feature a short description if the message alone is insufficient in conveying what has been changed. 

## Project

### Idea
The idea for the project is to create a 3D renderer and planet generator. The program will be written in the Go programming language and use OpenGL to send instructions to the GPU to draw the planets. The planets will be based on spheres and feature height differences to simulate mountains and terrain. The goal is to create realistic renders of planets, which requires both a thoughtfully constructed OpenGL renderer and well written algorithms to generate the planets. 

### Requirements 
Requirements structured according to MoSCoW.

#### Must have
* Functions for drawing 3D shapes to a program window. 

* Spherical planets with mountains and valleys.

* The planets must be procedurally generated, i.e. not simply imported from Blender. 

* A controllable camera to view the planets from different angles.

* Simple lighting from a single light source, the sun.

#### Should have
* A space skybox with stars.

* Craters.

* Different coloured surfaces for areas with different heights. 

* Transparent oceans.

#### Could have
* Mountains and trees cast shadows.

* Better looking oceans with waves and foam. 

* A realistic atmosphere with light refraction.

* Animated clouds. 

* A small solar system.

* Sliders to control the simulation.

#### Will not have
* Rings around planets.

* Terrain editing.

* Raytracing in surface lighting.

* A simulation of the entire observable universe.

### Feasibillity
Since we are three developers with at least some experience working with graphics programming, It should be feasable for us to accomplish our goal. The "must haves" seem entirely possible and even most "could haves" are not far from reasonable either. With enough dedication, we should be able to generate a nice looking solar system of neat planets.

### Areas of responsibility
Developers working on the project will choose issues to work on alone until completion, prioritizing issues scheduled earlier on the timeline. When done, they will merge their branch and choose another issue. No person working on the project will be strictly limited to any area or set of issues. Instead they are encouraged to share their skills and knowledge to help where they can. 
