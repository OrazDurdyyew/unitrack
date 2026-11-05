// Общая настройка тестов: i18n подключается ко всем компонентам, чтобы тексты брались из словаря.
import { config } from '@vue/test-utils';

import { i18n } from './providers/i18n';

config.global.plugins = [i18n];
