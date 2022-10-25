import React from 'react';
import { Badge, Flex, FlexItem, Text, TextContent, TextVariants } from '@patternfly/react-core';

function CidrBlockSideBar() {
    return (
        <Flex direction={{ default: 'column' }} flex={{ default: 'flex_1' }} className="pf-u-h-100">
            <Flex direction={{ default: 'row' }} className="pf-u-p-md pf-u-mb-0">
                <FlexItem>
                    <Badge style={{ backgroundColor: 'rgb(0,102,205)' }}>E</Badge>
                </FlexItem>
                <FlexItem>
                    <TextContent>
                        <Text component={TextVariants.h2} className="pf-u-font-size-xl">
                            External Entities
                        </Text>
                    </TextContent>
                </FlexItem>
            </Flex>
        </Flex>
    );
}

export default CidrBlockSideBar;