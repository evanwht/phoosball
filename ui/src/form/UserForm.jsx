import React from 'react';
import Form from 'react-bootstrap/Form';
import Row from 'react-bootstrap/Row';
import Button from 'react-bootstrap/Button';

const UserForm = props => {
    const {
        submit,
        submitButtonText,
        elements,
    } = props;

    let button = null;
    if (props.submitButtonText) {
        button = (
            <Row className="pr-1 pl-1 mb-4 mt-4">
                <Button variant="primary" type="submit" size="lg">{submitButtonText}</Button>
            </Row>
        );
    }

    return (
        <React.Fragment>
            <Form onSubmit={submit}>
                {elements()}
                {button}
            </Form>
        </React.Fragment >
    )
}

export default UserForm;