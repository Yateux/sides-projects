import React, {Component} from 'react';
import {BrowserRouter, Route, Switch} from "react-router-dom";


import TasksContainer from "./container/TasksContainer";
import HeaderContainer from "./container/HeaderContainer";
import CssBaseline from "@material-ui/core/CssBaseline/CssBaseline";
import {withStyles} from "@material-ui/core";
import NotFound from "./components/NotFound";
import Home from "./container/Home";

const styles = {
    root: {
        flexGrow: 1,
        zIndex: 1,
        overflow: "hidden",
        position: "relative",
        display: "flex"
    },
    content: {
        flexGrow: 1,
        maxWidth: '70vw',
        margin: 'auto',
        minWidth: 0
    }
};

class App extends Component {
    render() {
        const {classes} = this.props;
        return (
            <React.Fragment>
                <CssBaseline/>
                <BrowserRouter>
                    <React.Fragment>
                        <HeaderContainer/>
                        <div className={classes.root}>
                            <main className={classes.content}>
                                <Switch>
                                    <Route path="/tasks" component={TasksContainer}/>
                                    <Route path="/" exact component={Home}/>
                                    <Route path="*" component={NotFound}/>
                                </Switch>
                            </main>
                        </div>
                    </React.Fragment>
                </BrowserRouter>
            </React.Fragment>
        );
    }
}


export default withStyles(styles)(App);