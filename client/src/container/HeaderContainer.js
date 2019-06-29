import React from 'react';
import Toolbar from "@material-ui/core/Toolbar";
import AppBar from "@material-ui/core/AppBar";

import {withStyles} from "@material-ui/core";
import Link from "@material-ui/core/Link/Link";
import {Link as RouterLink} from 'react-router-dom';

const styles = {
    root: {
        marginBottom: 24
    },
    bar: {
        display: 'flex',
        justifyContent: 'space-between'
    },
    logout: {
        marginLeft: "50px"
    }
};

class HeaderContainer extends React.Component {

    // eslint-disable-next-line no-useless-constructor
    constructor(props) {
        super(props);
    }

    render() {
        const {classes} = this.props;

        return (
            <AppBar position="static" className={classes.root}>
                <Toolbar className={classes.bar}>
                    <Link component={RouterLink} to='/' color="inherit" variant="title">
                        Index
                    </Link>


                    <Link component={RouterLink} to='/tasks' color="inherit" variant="body1">
                        Tasks
                    </Link>


                </Toolbar>
            </AppBar>
        )
    }
}


export default withStyles(styles)(HeaderContainer);
